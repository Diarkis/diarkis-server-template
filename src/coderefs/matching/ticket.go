package matching

import (
	"sync"

	"github.com/Diarkis/diarkis/user"
)

/*
TicketParams parameter struct for issueTicket

	[IMPORTANT] AddProperties and SearchProperties are limited to have up to 2 properties.

Properties

	ProfileIDs                  - A list of profiles to add to and search against.
	AddProfileIDs               - Optional list of profiles to add. If this is used, ProfileIDs will be overridden for add.
	SearchProfileIDs            - Optional list of profiles to search. If this is used, ProfileIDs will be overridden for search.
	Tag                         - A string tag to group matchmaking data by the same tag.
	                              Data with different tag do NOT match even with the matching properties.
	                              Leave it with an empty string if no need.
	AddProperties               - Matchmaking properties (conditions) for add (being a host waiting)
	SearchProperties            - Matchmaking properties to used for search
	                              Each property may contain a range of property values i.e. []int{ 1, 2, 3, 4, 5 } etc.
	                              The number of properties allowed is 2,
	                              if you exceed the number of properties, two properties will be randomly chosen.
	ApplicationData             - May hold application data that maybe added from the application.
	                              ApplicationData must NOT be a struct or must NOT contain struct.
	MaxMembers                  - Maximum number of matchmaking users per matchmaking.
	                              When matched users reach this number, the matchmaking will complete as success
	TicketDuration              - Duration of the ticket to be valid in seconds.
	                              Minimum value for TicketDuration is 10 seconds.
	SearchInterval              - The interval for search in milliseconds
	TimeoutExtensionOnMatchJoin - Timeout extension in seconds to be added every time a new match joins. Leave it with 0 if no need.
	SearchTries                 - Number of empty search results to tolerate before giving up and moving on to hosting (add)
	EmptySearches               - If the number of empty search results reach EmptySearches, the ticket will forcefully change to add phase.
	                              If 0 is given, this feature will be disabled. Default is 0.
	HowMany                     - Matchmaking profile IDs to use for search and add
	                              Leave this empty if you do not need to repeat the operation set.

# Breaking down how a ticket works under the hood

A ticket has two phases. When you start a ticket,
it starts as a search phase where it actively searches other tickets that are in waiting phase.
Once a certain time passes and the search yields no matches,
ticket then switches to waiting phase where it waits for other searching tickets to find it.

Diagram below visually explains how these two phases of a ticket work and how some of the parameters affect these two phases.

	▶︎ SearchInterval 200ms
	▶︎ SearchTries    10
	▶︎ TicketDuration 4s

	▷ Search Phase: The ticket will search every 200ms 10 times
	▷ Wait Phase:   If the ticket does not find match, it will wait for the remainder of time until the ticket duration expires

	       200ms x 10 searches
	  ┌────── Search Phase ──────┐ ┌───────  Waiting Phase ───────┐
	┌──┬──┬──┬──┬──┬──┬──┬──┬──┬──┰────────────────────────────────┐
	│  │  │  │  │  │  │  │  │  │  ┃                                │
	└──┴──┴──┴──┴──┴──┴──┴──┴──┴──┸────────────────────────────────┘
	 Total 2 seconds of searching      Total 2 seconds of waiting
	│                                                              │
	└─────────── Ticket duration is 4 seconds in total ────────────┘

	This means that a ticket in search phase will only match with tickets in wait phase and vice versa.

# TicketParams Tip

It usually helps to have randomized values for SearchInterval, SearchTries, and TicketDuration.

This is because every ticket strictly follows search → wait flow.
Having every ticket with different search and wait durations will help tickets find other tickets.
*/
type TicketParams struct {
	ProfileIDs                  []string
	AddProfileIDs               []string
	SearchProfileIDs            []string
	Tag                         string
	AddProperties               map[string]int
	SearchProperties            map[string][]int
	ApplicationData             []byte
	MaxMembers                  uint8
	TicketDuration              uint8
	SearchInterval              uint16
	TimeoutExtensionOnMatchJoin uint8
	SearchTries                 uint8
	EmptySearches               uint8
	HowMany                     uint8
}

/*
Ticket represents a matchmaking ticket that manages a life cycle of issued ticket

	OnMatch                       - Raised when a remote user matches. By returning true, you may complete the ticket
	                                and raise OnComplete (OnComplete event is captured by matching.SetOnComplete callback)
	OnMatchedMemberJoined         - Raised when a matched member completes join.
	OnMatchedMemberJoinedAnnounce - Raised when a matched member completes join
	                                and returns ver, cmd, and message to be sent to all matched members.
	OnMatchedMemberLeaveAnnounce  - Raised when a matched member leaves and returns ver, cmd, and message to be sent to all matched members.
	OnMatchedMemberLeave          - Raised when a matched member user leave the match.
	OnTimeout                     - Raised when the ticket times out.
*/
type Ticket struct {
	OnMatch                       func(ticket *Ticket, userData *user.User, owner *user.User, roomID string, memberIDs []string) bool
	OnMatchedMemberJoined         func(ticket *Ticket, userData *user.User, owner *user.User, memberIDs []string)
	OnMatchedMemberLeave          func(ticket *Ticket, userData *user.User, owner *user.User, roomID string, memberIDs []string)
	OnTimeout                     func(userData *user.User)
	OnMatchedMemberJoinedAnnounce func(ticket *Ticket, userData, owner *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte)
	OnMatchedMemberLeaveAnnounce  func(ticket *Ticket, userData, owner *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte)
}

/*
TicketProperties represents both the owner of the matched ticket and the candidate to be matched.

It is primarily meant to be used for SetOnTicketAllowMatchIf callback.

	Owner            - Represents ticket owner's add and search properties (user that perform add/waiting).
	                   Owner add and search properties are pointers to the original properties
	                   and changing the values may influence the matchmaking.
	Candidates       - A map of  match candidate's add and search properties (user that performs searches)
	                   by candidate's UID as keys.
	                   Candidate add and search properties are pointers to the original properties
	                   and changing the values may influence the matchmaking.
*/
type TicketProperties struct{ sync.RWMutex }

/*
GetOwner returns the ticket owner's *TicketHolder instance

	[NOTE] Returned *TicketHolder is a copy of the actual candidate ticket holder.
*/
func (t *TicketProperties) GetOwner() *TicketHolder {
	return nil
}

/*
GetCandidateByUID returns the given UID's *TicketHolder instance.
If the second value is false, there is no candidate that matches the given UID.

	[NOTE] Returned *TicketHolder is a copy of the actual candidate ticket holder.
*/
func (t *TicketProperties) GetCandidateByUID(uid string) (*TicketHolder, bool) {
	return nil, false
}

/*
GetAllCandidates returns all candidates' *TicketHolder instances as a map.
The key of the map is UID.

	[NOTE] Returned map of *TicketHolder is a copy of the actual candidate ticket holders.
*/
func (t *TicketProperties) GetAllCandidates() map[string]*TicketHolder {
	return nil
}

/*
TicketHolder represents add and search properties of the ticket holder user.

	AddProperties    - Add properties of the ticket: Add properties are used to be found by other tickets.
	SearchProperties - Search properties of the ticket: Search properties are used to search for other tickets.
	ApplicationData  - May hold byte array encoded application data that maybe added from the application.
*/
type TicketHolder struct {
	AddProperties    map[string]int
	SearchProperties map[string][]int
	ApplicationData  []byte
}

/*
Client represents matchmaking ticket matched member user client.
*/
type Client struct {
	ID                  string
	SID                 string
	PublicAddress       string
	PrivateAddressBytes []byte
	UserData            map[string]interface{}
}

/*
ControlTicketParams [INTERNALLY USED ONLY]
*/
func ControlTicketParams(params *TicketParams) (int64, int64) {
	return 0, 0
}

/*
Start starts the life cycle of a ticket.

	[NOTE] Uses mutex lock internally.
*/
func (t *Ticket) Start() bool {
	return false
}

/*
IsTicketOwner returns true if the given user is the owner of its ticket

	[NOTE] Uses mutex lock internally.
*/
func IsTicketOwner(ticketType uint8, userData *user.User) bool {
	return false
}

/*
MarkTicketAsComplete finishes the ticket as complete immediately.

	[IMPORTANT] Only the owner user of the ticket may execute this function.

	[NOTE] Uses mutex lock internally.

Error Cases

	╒═══════════════════════╤═════════════════════════════════════════════════════════════════════════════╕
	│ Error                 │ Reason                                                                      │
	╞═══════════════════════╪═════════════════════════════════════════════════════════════════════════════╡
	│ Ticket not found      │ Either the ticket is not available or the given user does not own a ticket. │
	├───────────────────────┼─────────────────────────────────────────────────────────────────────────────┤
	│ User is not owner     │ User given is not the owner of the ticket.                                  │
	├───────────────────────┼─────────────────────────────────────────────────────────────────────────────┤
	│ Ticket room not found │ Internal room of the ticket is missing.                                     │
	│                       │ Most likely due of an internal bug causing the ticket data to be corrupt.   │
	╘═══════════════════════╧═════════════════════════════════════════════════════════════════════════════╛

Parameters

	ticketType - Ticket type.
	userData   - Owner user of the ticket to mark as complete.
*/
func MarkTicketAsComplete(ticketType uint8, userData *user.User) error {
	return nil
}

/*
MarkTicketAsCompleteWhenExpire marks the ticket as complete, but it will wait until the ticket expires.

	[IMPORTANT] Only the owner of the ticket may execute this function.

	[NOTE] Uses mutex lock internally.

This is useful when you need to have alternative conditions for matchmaking completion without having all expected members match.
*/
func MarkTicketAsCompleteWhenExpire(ticketType uint8, userData *user.User) bool {
	return false
}

/*
Stop interrupts the ticket and stops all matchmaking operations.

	[NOTE] Uses mutex lock internally.
*/
func (t *Ticket) Stop() bool {
	return false
}

/*
GetTicketType returns the ticket type of the *Ticket instance.
*/
func (t *Ticket) GetTicketType() uint8 {
	return 0
}

/*
GetRoomID returns the room ID of the ticket.
*/
func (t *Ticket) GetRoomID() string {
	return ""
}

/*
IsTicketFinished returns true if the ticket has finished its entire operations.
*/
func (t *Ticket) IsTicketFinished() bool {
	return false
}

/*
LeaveFromTicketMatchmaking makes the target user leave the matchmaking that the user has matched and joined.

	[NOTE] Uses mutex lock internally.
*/
func LeaveFromTicketMatchmaking(ticketType uint8, userData *user.User) bool {
	return false
}

/*
LeaveFromTicketMatchmakingWithCallback makes the target user leave the matchmaking that the user has matched and joined.

	[NOTE] Uses mutex lock internally.
*/
func LeaveFromTicketMatchmakingWithCallback(ticketType uint8, userData *user.User, cb func(err error)) {
}

/*
KickoutFromTicket forcefully removes a matched member from the matchmaking ticket.

	[IMPORTANT] Only the owner of the ticket may execute this operation.

	[NOTE] Uses mutex lock internally.

Error Cases

	╒═══════════════════════════╤═════════════════════════════════════════════════════════════════════════════╕
	│ Error                     │ Reason                                                                      │
	╞═══════════════════════════╪═════════════════════════════════════════════════════════════════════════════╡
	│ MatchMaker room not found │ The owner user is not in a matchmaking ticket room.                         │
	├───────────────────────────┼─────────────────────────────────────────────────────────────────────────────┤
	│ Must be the owner         │ Only the owner of the ticket may kick out matched members.                  │
	├───────────────────────────┼─────────────────────────────────────────────────────────────────────────────┤
	│ Target user not found     │ The target user to kick out is not a matched member.                        │
	╘═══════════════════════════╧═════════════════════════════════════════════════════════════════════════════╛
*/
func KickoutFromTicket(ticketType uint8, owner *user.User, targetUserID string, cb func(err error)) {
}

/*
SetOnIssueTicket assigns a callback to be invoked when the built-in command issue ticket is called.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked and expected to create and return TicketParams for a new ticket.

The callback is invoked for the owner user only.

Calling StartTicket triggers this callback and TicketParams that it returns will be used to create a new ticket.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] This callback must be assigned to the appropriate ticket type
	            in order for StartTicket(ticketType uint8, userData *user.User) to work properly.

	[IMPORTANT] userData maybe nil if the user disconnects from the server.
*/
func SetOnIssueTicket(ticketType uint8, cb func(userData *user.User) *TicketParams) bool {
	return false
}

/*
SetOnTicketAllowMatchIf assigns a callback to be invoked
before a match is made to add a custom logic to control if the match found should proceed forward to become an actual match or not..

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	callback   - Callback to be called func(ticketProps *TicketProperties, owner *user.User, candidate *user.User) bool

If the callback returns true, a match will be made.

The callback is invoked for the owner user only.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked while the ticket is holding the mutex lock,
	            In order to avoid mutex deadlock, you must not invoke functions that uses mutex lock or retrieve a mutex lock in the callback.
	[IMPORTANT] Candidate is the user that has been matched and attempting to join the matched member room.
	[IMPORTANT] owner and/or candidate maybe nil if owner and/or candidate disconnects from the server.
*/
func SetOnTicketAllowMatchIf(ticketType uint8, cb func(ticketProps *TicketProperties, owner, candidate *user.User) bool) bool {
	return false
}

/*
SetOnTicketMatch assigns a callback to be invoked when a new match is made.

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to control if the match completes the matchmaking or not.
	             Room ID passed to the callback is NOT Diarkis Room's ID, but it is the ID of the internal matchmaking room.
	             func(ticket *Ticket, matchedUserData *user.User, owner *user.User, roomID string, memberIDs []string) bool

This callback is meant to execute a custom logic for matchmaking completion on every match found.

Having the callback return true will automatically completes the matchmaking.

The callback is invoked for the owner user only.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] callback is invoked for the owner of the ticket only.
	[IMPORTANT] matchedUser is a copy of the actual user data of *user.User because the matched user maybe on a different server.
	[IMPORTANT] ownerUser and/or matchedUser maybe nil if ownerUser and/or matchedUser disconnects from the server.
*/
func SetOnTicketMatch(ticketType uint8, cb func(ticket *Ticket, matchedUser, ownerUser *user.User, roomID string, memberIDs []string) bool) bool {
	return false
}

/*
SetOnTicketMemberJoined assigns a callback to be invoked when a matched member joins.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked for the owner user only.
	[IMPORTANT] joinedUser is a copy of the actual user data of *user.User because the matched user maybe on a different server.
	[IMPORTANT] ownerUser and/or joinedUser maybe nil if ownerUser and/or joinedUser disconnects from the server.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a matched member joins the match.
*/
func SetOnTicketMemberJoined(ticketType uint8, cb func(ticket *Ticket, joinedUser, ownerUser *user.User, memberIDs []string)) bool {
	return false
}

/*
SetOnTicketMemberJoinedAnnounce assigns a callback to be invoked when a matched member joins
and returns ver, cmd, and message to be used as an announcement.

An announcement is sent to all matched users.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The event is invoked for the owner of the ticket only.

Parameters

	ticketType - Apply the callback to the given ticket type.
	cb         - Callback to be invoked when a new user successfully match and join.
	             ticket     - The matchmaking ticket.
	             joinedUser - The user that joined the matchmaking.
	                          [IMPORTANT] This is a copy data of the joined user and it is not the actual joined user data.
	             ownerUser  - The owner user of the matchmaking ticket.
	             memberIDs  - an array of matched member user IDs.

Example:

	added := matching.SetOnTicketMemberJoinedAnnounce(ticketType, func(ticket *matching.Ticket, joinedUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		// we will be sending a notification message to all matched users with the following:
		ver = uint8(2)                                  // the notification message will be sent with command ver 2
		cmd = uint16(1010)                              // the notification message will be sent with command ID 1010
		message = []byte(strings.Join(messageIDs, ",")) // the notification message will send a message with comma separated list of matched member user IDs.
		return ver, cmd, message
	})

	if !added {
	  // failed to assign the callback...
	}
*/
func SetOnTicketMemberJoinedAnnounce(ticketType uint8, cb func(ticket *Ticket, joinedUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte)) bool {
	return false
}

/*
SetOnTicketMemberLeave assigns a callback to be invoked when a matched member leaves.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked for the owner user only.
	[IMPORTANT] leftUser is a copy of the actual user data of *user.User because the matched user maybe on a different server.
	[IMPORTANT] ownerUser and/or leftUserData maybe nil if ownerUser and/or leftUserData disconnects from the server.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a matched member leaves the match.
*/
func SetOnTicketMemberLeave(ticketType uint8, cb func(ticket *Ticket, leftUserData, ownerUser *user.User, roomID string, memberIDs []string)) bool {
	return false
}

/*
SetOnTicketMemberLeaveAnnounce assigns a callback to be invoked when a matched member leaves
and returns ver, cmd, and message to be used as an announcement.

An announcement is sent to all matched users.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The event is invoked for the owner of the ticket only.

Parameters

	ticketType - Apply the callback to the given ticket type.
	cb         - Callback to be invoked when a matched user leaves.
	             ticket    - The matchmaking ticket.
	             leftUser  - The user that left the matchmaking.
	                         [IMPORTANT] This is a copy of the left user data and it is not the actual left user data.
	             ownerUser - The owner user of the matchmaking ticket.
	             memberIDs - an array of matched member user IDs.

Example:

	added := matching.SetOnTicketMemberLeaveAnnounce(ticketType, func(ticket *matching.Ticket, leftUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte) {
		// we will be sending a notification message to all matched users with the following:
		ver = uint8(2)                                  // the notification message will be sent with command ver 2
		cmd = uint16(1010)                              // the notification message will be sent with command ID 1010
		message = []byte(strings.Join(messageIDs, ",")) // the notification message will send a message with comma separated list of matched member user IDs.
		return ver, cmd, message
	})

	if !added {
	  // failed to assign the callback...
	}
*/
func SetOnTicketMemberLeaveAnnounce(ticketType uint8, cb func(ticket *Ticket, joinedUser, ownerUser *user.User, memberIDs []string) (ver uint8, cmd uint16, message []byte)) bool {
	return false
}

/*
SetOnTicketComplete assigns a callback to be invoked when an issued ticket successfully completes.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked only for the owner user of the ticket.

	[IMPORTANT] Owner maybe nil if owner disconnects from the server.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a matchmaking is complete.

The callback is invoked for the owner user of the ticket only.

The callback function must return a message byte array to be sent to all matched user clients.
*/
func SetOnTicketComplete(ticketType uint8, cb func(ticketProps *TicketProperties, owner *user.User) []byte) bool {
	return false
}

/*
SetOnTicketCanceled assigns a callback to be invoked when an issued ticket has been canceled.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked only for the owner user of the ticket.
	            For the non-owner users, OnMatchedTicketCancled will be invoked.

	[IMPORTANT] Owner maybe nil if owner is disconnected from the server.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a ticket is canceled.
*/
func SetOnTicketCanceled(ticketType uint8, cb func(owner *user.User)) bool {
	return false
}

/*
SetOnMatchedTicketCanceled assigns a callback to be invoked when a ticket that the user has joined is canceled.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked for the non-owner users of the ticket only.
	            For the ticket owner, OnTicketCancled will be invoked.
	[IMPORTANT] When the owner of the ticket either disconnects or re-connects to another server, the ticket will be canceled automatically.
	            Cancel event is raised event after the completion of the ticket when the owner of the ticket disconnects or re-connects..

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a ticket is canceled.
*/
func SetOnMatchedTicketCanceled(ticketType uint8, cb func(ownerID string, userData *user.User)) bool {
	return false
}

/*
SetOnTicketTimeout assigns a callback to be invoked when an issued ticket has timed out.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked for the owner user of the ticket only.
	[IMPORTANT] Owner maybe nil if owner disconnects from the server.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a ticket times out.
*/
func SetOnTicketTimeout(ticketType uint8, cb func(owner *user.User)) bool {
	return false
}

/*
SetOnMatchedTicketTimeout assigns a callback to be invoked when a ticket that the user has joined is timed out.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	[IMPORTANT] The callback is invoked for the non-owner users of the ticket only.
	            For the ticket owner, OnTicketTimeout will be invoked.

Parameters

	ticketType - Ticket type is used to group tickets. It will assign the callback to the given ticket type.
	cb         - Callback to be invoked when a ticket times out.
*/
func SetOnMatchedTicketTimeout(ticketType uint8, cb func(ownerID string, userData *user.User)) bool {
	return false
}

/*
StartTicket creates and starts a new ticket-based matchmaking with the given ticket type.

	[IMPORTANT] If the owner user (user that created the ticket) disconnects or re-connects, the ticket will be canceled automatically.
	            This also means that if you use Diarkis modules (Diarkis Field and Diarkis Room) that
	            may require the user to re-connect such as Room and Field,
	            it may cause the ticket to be canceled unexpected when the user re-connects.

	[IMPORTANT] Complex search properties with large number of properties with long range (lengthy array of values)
	            may have negative impact on server performance.

	[NOTE] Uses mutex lock internally.

Error Cases

	╒════════════════════════════════════════════════════╤════════════════════════════════════════════════════════════════════════════════╕
	│ Error                                              │ Reason                                                                         │
	╞════════════════════════════════════════════════════╪════════════════════════════════════════════════════════════════════════════════╡
	│ MatchMaker is not setup correctly                  │ SetOnIssueTicket callback for the ticketType must be assigned.                 │
	├────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ MatchMaker ticket cannot be issued more than once  │ The user is not allowed to issue more than one ticket of the same ticket type. │
	├────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Failed to create a ticket                          │ Given nil for *TicketParams.                                                   │
	├────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Failed to issue a new ticket                       │ The user is still in a previous matchmaking room of the same ticket type.      │
	╘════════════════════════════════════════════════════╧════════════════════════════════════════════════════════════════════════════════╛

Parameters

	ticketType - Ticket type is used to group tickets. You may issue multiple tickets of different ticket types.
	userData   - The user that issues and starts the ticket and becomes the owner of the ticket.

# How Ticket Works Internally

Ticket manages MatchMaker's Add and Search internally, so that you do not have to manage and balance their calls.

Ticket has two phases: "Search" and "Waiting"

Ticket starts out in "Search" phase and moves to "Waiting" phase.

The diagram below shows how a ticket operations internally:

	                                             ┌──────────────┐
	                                             │ Start Ticket │
	                                             └──────┬───────┘
	                                                    │
	                                                    ▼
	                                             ┌──────────────┐
	                                             │    Search    │
	                                             └──────┬───────┘
	                                                    │
	                                         ┌──────────┴──────────┐
	                                         │                     │
	    Ticket actively searches for         ▼                     ▼
	    other tickets that are       ┌───────────────┐        ┌────────┐ Ticket creates a waiting room for
	    in "Wait" phase to match     │ Found Matches │        │  Wait  │ the ticket and waits for
	                                 └──────┬────────┘        └──┬──┬──┘ other tickets to "Search" it and match
	╭────────────────────────────────────╮  │                    │  │
	│ Ticket found other tickets         │  └────────┐    ┌──────┘  └──────────┐
	│ to join and matched ticket         ├──────────▷│    │◁──────────┐        │
	│ met the required number of users   │           │    │           │        │
	╰────────────────────────────────────╯           ▼    ▼           │        ▼
	                                      ╔═══════════════════════╗   │   ┌─────────┐ If required number of users are not met
	                                      ║ Matchmaking Completed ║   │   │ Timeout │ and ticket duration expires,
	                                      ╚═══════════════════════╝   │   └─────────┘ the ticket times out. In order to continue,
	                                                       ┌──────────┘               The user must issue and start a new ticket
	                                                       │
	    ╭──────────────────────────────────────────────────┴─╮
	    │ Other tickets searched and found the ticket        │
	    │ and the required number of users have been matched │
	    ╰────────────────────────────────────────────────────╯

Diagram below visually explains how these two phases of a ticket work and how some of the parameters affect these two phases.

	▶︎ SearchInterval 200ms
	▶︎ SearchTries    10
	▶︎ TicketDuration 4s

	▷ Search Phase: The ticket will search every 200ms 10 times
	▷ Wait Phase:   If the ticket does not find match, it will wait for the remainder of time until the ticket duration expires

	       200ms x 10 searches
	  ┌────── Search Phase ──────┐ ┌───────  Waiting Phase ───────┐
	┌──┬──┬──┬──┬──┬──┬──┬──┬──┬──┰────────────────────────────────┐
	│  │  │  │  │  │  │  │  │  │  ┃                                │
	└──┴──┴──┴──┴──┴──┴──┴──┴──┴──┸────────────────────────────────┘
	 Total 2 seconds of searching      Total 2 seconds of waiting
	│                                                              │
	└─────────── Ticket duration is 4 seconds in total ────────────┘

	This means that a ticket in search phase will only match with tickets in wait phase and vice versa.

# Matchmaking result notifications

The table below explains the notifications that the client receives from the server.

	┌───────────────────┬──────────────────┬─────────────────┬───────────────────────────────────────────────────────────────────────────┐
	│ Notification Type │ Push Command Ver │ Push Command ID │ Description                                                               │
	╞═══════════════════╪══════════════════╪═════════════════╪═══════════════════════════════════════════════════════════════════════════╡
	│ Success           │ 1                │ 220             │ Matchmaking ticket has been successfully completed                        │
	│                   │                  │                 │ and all matched user clients receive this server push.                    │
	├───────────────────┼──────────────────┼─────────────────┼───────────────────────────────────────────────────────────────────────────┤
	│ Timeout           │ 1                │ 219             │ Matchmaking ticket has failed and it has been discarded.                  │
	│                   │                  │                 │ All user clients that matched receives this server push.                  │
	├───────────────────┼──────────────────┼─────────────────┼───────────────────────────────────────────────────────────────────────────┤
	│ Cancel            │ 1                │ 222             │ Matchmaking ticket has been canceled                                      │
	│                   │                  │                 │ and all user clients that matched receives this server push.              │
	├───────────────────┼──────────────────┼─────────────────┼───────────────────────────────────────────────────────────────────────────┤
	│ Broadcast         │ 1                │ 224             │ Matchmaking ticket sends a broadcast message to all matched user clients. │
	└───────────────────┴──────────────────┴─────────────────┴───────────────────────────────────────────────────────────────────────────┘

Calling StartTicket raises The callback assigned by SetOnIssueTicket and a new ticket will be created using the given ticket parameters.

Example with SetOnIssueTicket callback

	// the callback will be invoked by matching.StartTicket
	matching.SetOnIssueTicket(sampleTicketType, func(userData *user.User) *matching.TicketParams {
		return &matching.TicketParams{
			ProfileIDs:       []string{"RankMatch"},
			MaxMembers:       2,
			SearchInterval:   300, // 300 milliseconds interval of search
			SearchTries:      4,   // allow 4 consecutive empty search results up to 4 times before moving on to wait
			TicketDuration:   20,  // ticket lasts for 20 seconds
			HowMany:          20,  // up to 20 search results
			Tag:              "",  // if we want to group tickets using tag add the string value here
			AddProperties:    &map[string]int{ "Rank": 3 }, // wait for other users to find me and my rank is 3
			SearchProperties: &map[string][]int{ "Rank": &[]int{ 1, 2, 3, 4, 5 } }, // search for other users within the range of 1 to 5 and property is "Rank"
		}
	})

	err := matching.StartTicket(sampleTicketType, userData)

	if err != nil {
		// error...
	}

▶︎ SearchPropeties

ServerProperties dictate conditions for searches that is performed internally.

	[IMPORTANT] The number of allowed range properties is limited to two.
	            When you have multiple elements in a search property, it is considered as a range property.

	[IMPORTANT] SearchProperties operates AND operations. It means that with multiple search properties,
	            the search must satisfy all search properties in order to match.

Diagram below shows how each search property operates:

	     ┏━━━┓
	     ┃ 8 ┃ 8 with ±3 would fall into 5 ~ 11 and that means it matches with items in the bucket of 0 to 10 and 11 to 20
	     ┗━┳━┛
	     ┏━┻━━━━┓
	│    ▼    │ ▼       │         │
	│  0 ~ 10 │ 11 ~ 20 │ 21 ~ 30 │
	└─────────┴─────────┴─────────┘

Example:

	searchProps["level"]     = []int{ 1, 2, 3, 4, 5 } // range property
	searchProps["rank"]      = []int{ 1, 2, 3 } // range property
	searchProps["matchType"] = []int{ 1 } // regular property
	searchProps["league"]    = []int{ 10 } // regular property

Order of range search property's search:

The range search properties will look for matches in the order of the array.
It means that if []int{ 1, 2, 3, 4, 5 } is given, it will start from 1 and continue up to 5 until it finds it matches.

▶︎ Callbacks

Every callback is assigned to given ticket type and invoked based on assigned ticket type.
*/
func StartTicket(ticketType uint8, userData *user.User) error {
	return nil
}

/*
StartTicketBackfill starts "backfill" on a ticket that has already been completed.

Backfill allows other users to match and join "after" the completion of the ticket.

When ticket is completed, the ticket itself will be deleted, but the matched users remain in the ticket room.

The owner user that wishes to start backfill must be a member of this ticket room.

	[IMPORTANT] If the ticket is full, no users may match and join even with backfill started.
	[IMPORTANT] In order to stop backfill, you must invoke StopTicketBackfill.

	[NOTE] Uses mutex lock internally.

Error Cases

	┌────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────┐
	│ Error                                              │ Reason                                                                         │
	╞════════════════════════════════════════════════════╪════════════════════════════════════════════════════════════════════════════════╡
	│ Completed ticket room must be available            │ The ticket and its internal room have been discarded.                          │
	│ Ticket already exists                              │ Either the ticket has not been completed or backfill has already been started. │
	│ MatchMaker is not setup correctly                  │ SetOnIssueTicket callback for the ticketType must be assigned.                 │
	│ MatchMaker ticket cannot be issued more than once  │ The user is not allowed to issue more than one ticket of the same ticket type. │
	│ Failed to create a ticket                          │ Given nil for *TicketParams.                                                   │
	│ Failed to issue a new ticket                       │ The user is still in a previous matchmaking room of the same ticket type.      │
	└────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────┘

Parameters

	ticketType - Type of the completed ticket to start backfill.
	owner      - Ticket's owner user.
*/
func StartTicketBackfill(ticketType uint8, owner *user.User) error {
	return nil
}

/*
StopTicketBackfill stops backfill ticket.

	[NOTE] Uses mutex lock internally.

	Error Cases

	┌────────────────────────────────┬─────────────────────────────────────────────────────────────────────────┐
	│ Error                          │ Reason                                                                  │
	╞════════════════════════════════╪═════════════════════════════════════════════════════════════════════════╡
	│ Backfill ticket not found      │ Backfill ticket to stop does not exist.                                 │
	│ Backfill ticket room not found │ Backfill ticket room has been discarded. (All users have disconnected.) │
	└────────────────────────────────┴─────────────────────────────────────────────────────────────────────────┘

Parameters

	ticketType - Backfill ticket type to stop backfill.
	owner      - Backfill ticket owner user.
*/
func StopTicketBackfill(ticketType uint8, owner *user.User) error {
	return nil
}

/*
HasTicket returns true if the user given has a matchmaking ticket of the given type in progress.
*/
func HasTicket(ticketType uint8, userData *user.User) bool {
	return false
}

/*
CancelTicket stops ticket-based matchmaking started by StartTicket.

Canceling a ticket will disband the matchmaking and remove already-matched users from it.

	[IMPORTANT] Only the owner user of the ticket are allowed to cancel the ticket.
	[IMPORTANT] CancelTicket will fail after the completion of the ticket.

	[NOTE] Uses mutex lock internally.

Error Cases

	┌───────────────────────┬───────────────────────────────────────────────────────────┐
	│ Error                 │ Reason                                                    │
	╞═══════════════════════╪═══════════════════════════════════════════════════════════╡
	│ Ticket not found      │ The ticket to cancel is not available.                    │
	│ Ticket failed to stop │ The ticket is not active (already completed or canceled). │
	└───────────────────────┴───────────────────────────────────────────────────────────┘

Parameters

	ticketType - Ticket type of the ticket to cancel.
	userData   - The owner user of the ticket to cancel.
*/
func CancelTicket(ticketType uint8, userData *user.User) error {
	return nil
}

/*
FindTicket returns the valid matchmaking ticket that the user has.

	[IMPORTANT] This function works with the owner of the ticket only.
	[IMPORTANT] IF the user given does not own a ticket, it returns nil.

	[NOTE] Uses mutex lock internally.

	ticketType - MatchMaker ticket type.
	userData   - The owner user of the ticket.
*/
func FindTicket(ticketType uint8, userData *user.User) *Ticket {
	return nil
}
