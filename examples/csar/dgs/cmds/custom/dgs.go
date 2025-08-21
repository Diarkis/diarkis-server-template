// © 2019-2025 Diarkis Inc. All rights reserved.

package customcmds

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/Diarkis/diarkis/derror"
	"github.com/Diarkis/diarkis/dgs"
	"github.com/Diarkis/diarkis/dgs/types"
	"github.com/Diarkis/diarkis/diarkisexec"
	"github.com/Diarkis/diarkis/room"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"

	pufferDgs "github.com/Diarkis/diarkis-server-template/examples/csar/dgs/puffer/go/dgs"
)

// DGS instance create response.
// The status possible values are:
//   - 0 success (created)
//   - 1 DGS is already starting. A previous call to instance create have been done and it is still ongoing.
//       The caller is part of the DGS creation members request.
//   - 2 DGS is already started. Client should either call backfill or the owner should call clear association.
//   - 3 DGS is starting and the caller is not part of the DGS creation members request.
//       Caller will not receive a credentials push but will receive InstanceAllocatedPush
//       with IsMember=false when the creation will be finished. Caller should call backfill
//       once this push event is received.

// DGS backfill response.
// The status possible values are:
//   - 0 success
//   - 1 DGS is starting. A previous call to instance create have been done and it is still ongoing.
//   - 2 DGS instance not created. The owner of the session must first call create DGS instance.
//   - 3 DGS is starting and the caller is not part of the DGS creation members request.
//       Caller will not receive a credentials push but will receive InstanceAllocatedPush
//       with IsMember=false when the creation will be finished. Caller should call backfill again
//       once this push event is received.

// room property DGS creation state
const (
	dgsStateCreating string = "creating"
	dgsStateCreated  string = "created"
)

// DGS create response statuses.
const (
	// DGS has been successfully created
	dgsCreationStatusCreated uint8 = 0
	// DGS creation is ongoing, and the caller is contained in the DGS creation request member.
	// Caller should receive a push with the credentials when finished.
	dgsCreationStatusCreatingAsMember uint8 = 1
	// DGS creation is already done. In order to create a new DGS instance
	// the current one must be terminated.
	dgsCreationStatusAlreadyCreated uint8 = 2
	// DGS creation is ongoing, and the caller is not contained in the DGS creation request member.
	// Caller should receive a push when the creation is finished but no push with credentials.
	dgsCreationStatusCreatingNotAsMember uint8 = 3
)

// DGS backfill response statuses.
const (
	// DGS backfill success. The caller will soon receive a push with its credentials.
	dgsBackfillStatusSuccess uint8 = 0
	// DGS creation is ongoing, and the caller is contained in the DGS creation request member.
	// Caller should receive a push with the credentials when finished.
	dgsBackfillStatusDGSIsStartingAsMember uint8 = 1
	// No DGS is associated to the room. The owner must call create instance first.
	dgsBackfillStatusDGSNotCreated uint8 = 2
	// DGS creation is ongoing, and the caller is not contained in the DGS creation request member.
	// Caller should receive a push when the creation is finished but no push with credentials.
	dgsBackfillStatusDGSIsStartingNotAsMember uint8 = 3
)

const roomDGSStatePropertyName = "__dgsState"

// dgsAllocationTimeout how long the call to dgs.InstanceCreate
// is allowed to take.
//
// FIXME
// The timeout below should be long enough to allow the real DGS
// to handle the start session RPC. If the DGS needs to load
// some asset while processing the request, this should be take
// into account.
const dgsAllocationTimeout = time.Second * 20

// dgsCreateSessionFromRoom Instantiate a DGS using the room members.
// The caller MUST be the room owner.
// Any error that is not caller is not the owner of the room will trigger
// InstanceAllocatedPush with success=false ONLY and ONLY is status is not diarkisexec.CommandResponseOK.
func dgsCreateSessionFromRoom(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Info("dgsCreateSessionFromRoom: user %s payload: %x", userData.ID, payload)
	reqProto := pufferDgs.NewInstanceCreateRequest()
	if err := reqProto.Unpack(payload); err != nil {
		userData.ServerRespond(derror.ErrData("payload unpack error", derror.InvalidParameter(0)), ver, cmd, diarkisexec.CommandResponseBad, true)
		next(util.NewError("payload unpack error %v -> UserID:%s", err, userData.ID))
		return
	}

	roomID := room.GetRoomID(userData)
	if roomID == "" {
		// The user is not in a room.
		userData.ServerRespond(derror.ErrData("User not in the room", derror.ConditionFailed(0)), ver, cmd, diarkisexec.CommandResponseBad, true)
		next(util.NewError("User not in the room -> RoomID:%s UserID:%s", roomID, userData.ID))
		return
	}

	memberIDs, ownerID, _ := room.GetMemberIDsAndOwner(roomID)
	if ownerID != userData.ID {
		// Only the owner of the room can start the DGS
		userData.ServerRespond(derror.ErrData("User is not the room owner", derror.ConditionFailed(0)), ver, cmd, diarkisexec.CommandResponseBad, true)
		next(util.NewError("User is not the room owner -> RoomID:%s UserID:%s", roomID, userData.ID))
		return
	}

	// currentState pointer will be the same as creationState
	// only if the update was successful.
	var currentState *pufferDgs.RoomDGSState

	creationState := pufferDgs.NewRoomDGSState()
	creationState.LastUpdateTime = time.Now().UnixMicro()
	creationState.State = dgsStateCreating
	creationState.Members = memberIDs

	currentState, updatePropertiesErr := updateRoomDGSState(roomID, creationState)

	if updatePropertiesErr != nil {
		// An error occured.
		logger.Errorf("updatePropertiesErr", "Error", updatePropertiesErr)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, memberIDs, ownerID, roomID)
		next(util.NewError("failed to unpack room dgs state -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, updatePropertiesErr))
		return
	}

	if currentState == creationState {
		// We successfully updated the room property.
		// Move on to the DGS creation phase.
		// Run createDGS in a goroutine because it can take an arbitrary amount of time
		// to finish.
		next(nil)

		proto := pufferDgs.NewMeshInstanceCreateRequest()
		proto.Metadata = nil // application custom data to send to the DGS server
		proto.MemberIDs = memberIDs

		go createDGS(ver, cmd, roomID, ownerID, userData, proto)
		return
	}

	if currentState.State == dgsStateCreating {
		// Already in creation phase.
		logger.Infof("Cannot allocate DGS", "Room", roomID, "User", userData, "DGSStatus", currentState.State)
		status := dgsCreationStatusCreatingAsMember
		if !slices.Contains(currentState.Members, userData.ID) {
			status = dgsCreationStatusCreatingNotAsMember
		}
		response := pufferDgs.InstanceCreateResponse{
			Status: status,
		}
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	}

	if currentState.State != dgsStateCreated {
		// invalid state
		// return an error
		logger.Errorf("Cannot allocate DGS. Invalid state", "Room", roomID, "User", userData, "DGSStatus", currentState.State)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, memberIDs, ownerID, roomID)
		next(util.NewError("cannot allocate DGS. Invalid state -> RoomID:%s UserID:%s", roomID, userData.ID))
		return
	}

	// Check if the currently associated DGS server is still running.
	// If it is still running returns an error,
	// otherwise move in creation phase.

	logger.Infof("current state is created. Check is the dgs identifier is still valid")

	// If the DGS is already created, check if the associated DGS is still alive.
	logger.Info("current state identifier %x", currentState.DgsIdentifier)
	// Here we use a background context because the RPC is supposed to be really fast.
	err := dgs.InstanceValidateToken(context.Background(), currentState.DgsIdentifier)
	if err == nil {
		response := pufferDgs.InstanceCreateResponse{
			Status: dgsCreationStatusAlreadyCreated,
		}
		logger.Infof("Cannot allocate DGS", "Room", roomID, "User", userData, "DGSStatus", currentState.State)
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	}

	logger.Info("dgs.InstanceValidateToken", "Error", err)
	if !errors.Is(err, types.ErrCodeInvalidDGSIdentifier) {
		// internal error
		logger.Errorf("Cannot allocate DGS", "Room", roomID, "User", userData, "Error", err)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, memberIDs, ownerID, roomID)
		next(util.NewError("Cannot allocate DGS -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, err))
		return
	}

	// Before updating the room DGS state, we need to clean the property,
	// otherwise the call to updateRoomDGSState below will fail.
	dgsResetRoomDGSState(roomID)

	// Try to update the room property and take ownership of the creation.
	currentState, updatePropertiesErr = updateRoomDGSState(roomID, creationState)

	if updatePropertiesErr != nil {
		// An error occured.
		logger.Errorf("updatePropertiesErr", "Error", updatePropertiesErr)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, memberIDs, ownerID, roomID)
		next(util.NewError("failed to unpack room dgs state -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, updatePropertiesErr))
		return
	}

	if currentState == creationState {
		// We successfully updated the room property.
		// Move on to the DGS creation phase.
		// Run createDGS in a goroutine because it can take an arbitrary amount of time
		// to finish.
		next(nil)

		proto := pufferDgs.NewMeshInstanceCreateRequest()
		proto.Metadata = nil // application custom data to send to the DGS server
		proto.MemberIDs = memberIDs

		go createDGS(ver, cmd, roomID, ownerID, userData, proto)
		return
	}

	if currentState.State == dgsStateCreating {
		// Already in creation phase.
		logger.Infof("Cannot allocate DGS", "Room", roomID, "User", userData, "DGSStatus", currentState.State)
		status := dgsCreationStatusCreatingAsMember
		if !slices.Contains(currentState.Members, userData.ID) {
			status = dgsCreationStatusCreatingNotAsMember
		}
		response := pufferDgs.InstanceCreateResponse{
			Status: status,
		}
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	}

	if currentState.State == dgsStateCreated {
		// Already in creation phase.
		logger.Infof("Cannot allocate DGS", "Room", roomID, "User", userData, "DGSStatus", currentState.State)
		response := pufferDgs.InstanceCreateResponse{
			Status: dgsCreationStatusAlreadyCreated,
		}
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	}

	// Consider all other cases as an error occured.
	userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
	sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, memberIDs, ownerID, roomID)
	next(util.NewError("failed to unpack room dgs state -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, updatePropertiesErr))
}

// updateRoomDGSState update the room DGS property with the newState if possible.
// If the update is not possible, returns the old state.
func updateRoomDGSState(roomID string, newState *pufferDgs.RoomDGSState) (currentState *pufferDgs.RoomDGSState, updatePropertiesErr error) {
	room.UpdateProperties(roomID, func(properties map[string]any) bool {
		stateAny, ok := properties[roomDGSStatePropertyName]
		if !ok {
			properties[roomDGSStatePropertyName] = newState.Pack()
			currentState = newState
			return true
		}
		b, ok := util.ToBytes(stateAny)
		if !ok {
			// weird
			updatePropertiesErr = fmt.Errorf("room property %s type is invalid. got %T, want []byte", roomDGSStatePropertyName, stateAny)
			return false
		}
		var tmp pufferDgs.RoomDGSState
		if err := tmp.Unpack(b); err != nil {
			updatePropertiesErr = err
			return false
		}
		currentState = &tmp

		return false
	})

	return
}

// createDGS Create the DGS with the given request.
// On error reset the room's DGS state and push a failure event to the members.
func createDGS(ver uint8, cmd uint16, roomID, ownerID string, userData *user.User, data *pufferDgs.MeshInstanceCreateRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), dgsAllocationTimeout)
	defer cancel()

	dgsResponse, dgsIdentifier, err := dgs.InstanceCreate(ctx, data.Pack(), nil, nil)
	if err != nil {
		logger.Errorf("Failed to allocate DGS", "Room", roomID, "User", userData, "Error", err)
		dgsResetRoomDGSState(roomID)
		userData.ServerRespond(derror.ErrData("Failed to allocate DGS.", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, data.MemberIDs, ownerID, roomID)
		sendDGSAllocationFailurePush(ver, cmd, data.MemberIDs, ownerID)
		return
	}

	// Parse the DGS response using puffer
	// and create one PushInstanceCreate per user to send their credentials.
	meshResponse := pufferDgs.NewMeshInstanceCreateResponse()
	if err := meshResponse.Unpack(dgsResponse); err != nil {
		logger.Errorf("Failed to unpack DGS create response", "Room", roomID, "User", userData, "Error", err)
		dgsResetRoomDGSState(roomID)
		userData.ServerRespond(derror.ErrData("Failed to allocate DGS", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, false, data.MemberIDs, ownerID, roomID)
		sendDGSAllocationFailurePush(ver, cmd, data.MemberIDs, ownerID)
		return
	}

	state := pufferDgs.NewRoomDGSState()
	state.LastUpdateTime = time.Now().UnixMicro()
	state.State = dgsStateCreated
	state.DgsIdentifier = dgsIdentifier
	state.Members = data.MemberIDs

	logger.Info("current state new identifier %x", state.DgsIdentifier)

	// Store the DGS identifier in the room property so we can do RPC later.
	// Here we assume there is no possibility of concurrent DGS creation,
	// it should be fine to override the DGS property.
	room.UpdateProperties(roomID, func(properties map[string]interface{}) bool {
		properties[roomDGSStatePropertyName] = state.Pack()
		return true
	})

	response := pufferDgs.InstanceCreateResponse{
		Status: dgsCreationStatusCreated,
	}
	userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)

	sendDGSInstanceAllocatedPush(pufferDgs.InstanceAllocatedPushVer, pufferDgs.InstanceAllocatedPushCmd, true, data.MemberIDs, ownerID, roomID)

	for _, credentials := range meshResponse.Credentials {
		memberUser := user.GetUserByUID(credentials.UID)
		if memberUser == nil {
			logger.Warnf("Room member not found", "UID", credentials.UID)
			continue
		}

		proto := pufferDgs.NewInstanceCreatePush()
		proto.Success = true
		proto.SID = credentials.SID
		proto.IV = credentials.IV
		proto.Key = credentials.Key
		proto.MacKey = credentials.MacKey
		proto.ServerHost = credentials.ServerHost
		proto.ServerPort = credentials.ServerPort
		proto.ServerType = credentials.ServerType

		logger.Info("createDGS: send push to %s. %s", memberUser.ID, proto.String())
		memberUser.ServerPush(ver, cmd, proto.Pack(), true)
	}
}

func forEachMember(memberIDs []string, ownerID string, f func(*user.User)) {
	for _, memberID := range memberIDs {
		memberUser := user.GetUserByUID(memberID)
		if memberUser == nil {
			logger.Warnf("Room member not found", "UID", memberID)
			continue
		}
		// In case of error the owner will first receive the command error response.
		// No need to send a push event in that case.
		if memberUser.ID == ownerID {
			continue
		}

		f(memberUser)
	}
}

// sendDGSInstanceAllocatedPush Sends a push to all the member of the room.
func sendDGSInstanceAllocatedPush(ver uint8, cmd uint16, success bool, dgsMemberIDs []string, ownerID, roomID string) {
	proto := pufferDgs.NewInstanceAllocatedPush()
	proto.Success = success

	forEachMember(room.GetMemberIDs(roomID), ownerID, func(u *user.User) {
		proto.IsMember = slices.Contains(dgsMemberIDs, u.ID)
		logger.Info("sendDGSInstanceAllocatedPush: send push to %s. %s", u.ID, proto.String())
		u.ServerPush(ver, cmd, proto.Pack(), true)
	})
}

// sendDGSAllocationFailurePush Send a push to all member with Success set to false.
func sendDGSAllocationFailurePush(ver uint8, cmd uint16, memberIDs []string, ownerID string) {
	// create response once
	proto := pufferDgs.NewInstanceCreatePush()
	proto.Success = false
	proto.SID = nil
	proto.IV = nil
	proto.Key = nil
	proto.MacKey = nil
	proto.ServerHost = ""
	proto.ServerPort = 0
	proto.ServerType = ""

	forEachMember(memberIDs, ownerID, func(u *user.User) {
		logger.Info("sendDGSAllocationFailurePush: send push to %s", u.ID)
		u.ServerPush(ver, cmd, proto.Pack(), true)
	})
}

// dgsSessionBackfillFromRoom Allow a client to join an already existing DGS.
func dgsSessionBackfillFromRoom(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Info("dgsSessionBackfillFromRoom: user %s payload: %x", userData.ID, payload)
	reqProto := pufferDgs.NewInstanceBackfillRequest()
	if err := reqProto.Unpack(payload); err != nil {
		userData.ServerRespond(derror.ErrData("payload unpack error", derror.InvalidParameter(0)), ver, cmd, diarkisexec.CommandResponseBad, true)
		next(util.NewError("payload unpack error %v -> UserID:%s", err, userData.ID))
		return
	}

	roomID := room.GetRoomID(userData)
	if roomID == "" {
		// The user is not in a room
		userData.ServerRespond(derror.ErrData("User not in the room", derror.ConditionFailed(0)), ver, cmd, diarkisexec.CommandResponseBad, true)
		next(util.NewError("User not in the room -> RoomID:%s UserID:%s", roomID, userData.ID))
		return
	}

	// Retrieve the DGS identifier from the room property.
	properties := room.GetProperties(roomID)
	stateAny, ok := properties[roomDGSStatePropertyName]
	if !ok {
		logger.Debugf("Cannot backfill DGS", "Room", roomID, "User", userData, "Reason", "DGS not yet created")
		response := pufferDgs.InstanceBackfillResponse{}
		response.Status = dgsBackfillStatusDGSNotCreated
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	}

	b, ok := util.ToBytes(stateAny)
	if !ok {
		// weird
		err := fmt.Errorf("room property %s type is invalid. got %T, want []byte", roomDGSStatePropertyName, stateAny)
		logger.Errorf("Cannot backfill DGS", "Room", roomID, "User", userData, "Error", err)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		next(util.NewError("cannot backfill DGS -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, err))
		return
	}
	var state pufferDgs.RoomDGSState
	if err := state.Unpack(b); err != nil {
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		next(util.NewError("failed to unpack room dgs state -> RoomID:%s UserID:%s Error:%v", roomID, userData.ID, err))
		return
	}

	switch state.State {
	case dgsStateCreating:
		logger.Debugf("Cannot backfill DGS", "Room", roomID, "User", userData, "Reason", "DGS creation ongoing")
		response := pufferDgs.InstanceBackfillResponse{}
		if slices.Contains(state.Members, userData.ID) {
			response.Status = dgsBackfillStatusDGSIsStartingNotAsMember
		} else {
			response.Status = dgsBackfillStatusDGSIsStartingAsMember
		}
		userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
		next(nil)
		return
	case dgsStateCreated:
		// OK
	default:
		// Should not happen
		logger.Errorf("Cannot backfill DGS because room state is invalid", "Room", roomID, "User", userData, "State", state.State)
		userData.ServerRespond(derror.ErrData("room state is invalid", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		next(util.NewError("room state is invalid -> RoomID:%s UserID:%s", roomID, userData.ID))
		return
	}

	// Run DGS RPC in a goroutine because it can take an arbitrary amount of time
	// to finish.
	go next(nil)

	proto := pufferDgs.NewMeshInstanceCreateRequest()
	proto.Metadata = nil // application custom data to send to the DGS server
	proto.MemberIDs = []string{userData.ID}

	rpc := pufferDgs.MeshInstanceRPCRequest{}
	rpc.Data = proto.Pack()
	rpc.Type = "backfill"

	ctx, cancel := context.WithTimeout(context.Background(), dgsAllocationTimeout)
	defer cancel()
	response, err := dgs.InstanceRPC(ctx, rpc.Pack(), state.DgsIdentifier)
	if err != nil {
		logger.Errorf("failed to contact DGS", "Room", roomID, "User", userData, "Error", err)
		if errors.Is(err, types.ErrUnknownError) {
			userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
			return
		}
		if errors.Is(err, types.ErrCodeInvalidDGSIdentifier) {
			// DGS identifier missmatch means the session has been destroyed
			logger.Debugf("Cannot backfill DGS", "Room", roomID, "User", userData, "Reason", "DGS identifier missmatch")
			response := pufferDgs.InstanceBackfillResponse{}
			// From the point of view of the client it is similar to DGS not been created.
			response.Status = dgsBackfillStatusDGSNotCreated
			userData.ServerRespond(response.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
			return
		}

		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		return
	}

	credentials := pufferDgs.MeshInstanceCreateResponseUserCredentials{}
	if err := credentials.Unpack(response); err != nil {
		logger.Errorf("failed to parse DGS response", "Room", roomID, "User", userData, "Error", err)
		userData.ServerRespond(derror.ErrData("internal error", derror.Internal(0)), ver, cmd, diarkisexec.CommandResponseErr, true)
		return
	}

	{
		// Respond to the user
		proto := pufferDgs.InstanceBackfillResponse{}
		proto.Status = dgsBackfillStatusSuccess
		userData.ServerRespond(proto.Pack(), ver, cmd, diarkisexec.CommandResponseOK, true)
	}

	// Push the credentials to the user
	pushProto := pufferDgs.NewInstanceCreatePush()
	pushProto.Success = true
	pushProto.SID = credentials.SID
	pushProto.IV = credentials.IV
	pushProto.Key = credentials.Key
	pushProto.MacKey = credentials.MacKey
	pushProto.ServerHost = credentials.ServerHost
	pushProto.ServerPort = credentials.ServerPort
	pushProto.ServerType = credentials.ServerType

	logger.Info("dgsSessionBackfillFromRoom: send push to %s. %s", userData.ID, pushProto.String())
	userData.ServerPush(ver, cmd, pushProto.Pack(), true)
}

func dgsResetRoomDGSState(roomID string) {
	room.UpdateProperties(roomID, func(properties map[string]interface{}) bool {
		delete(properties, roomDGSStatePropertyName)
		return true
	})
}
