package lod

import (
	"testing"
	"time"
)

func TestProcessSenderReceiverPair(t *testing.T) {
	nearbyDistance := int32(100)
	farDistance := int32(500)
	nearbyInterval := 100 * time.Millisecond
	farInterval := 500 * time.Millisecond
	m := NewManager(1, 1, nearbyInterval, farInterval, nearbyDistance, farDistance)

	tests := []struct {
		name           string
		senderID       string
		receiverID     string
		senderPos      [2]int32
		receiverPos    [2]int32
		changed        bool
		lastSent       time.Duration
		expectedResult string
	}{
		{
			name:           "Self (Sender == Receiver)",
			senderID:       "user1",
			receiverID:     "user1",
			senderPos:      [2]int32{100, 100},
			receiverPos:    [2]int32{100, 100}, // same user
			expectedResult: "",
		},
		{
			name:           "Same position",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{100, 100},
			receiverPos:    [2]int32{100, 100}, // same position
			lastSent:       time.Duration(farInterval) + 10*time.Millisecond,
			changed:        false,
			expectedResult: "",
		},
		{
			name:           "Farther than maxDistanceForFar",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{100, 100},
			receiverPos:    [2]int32{600, 0}, // Distance 600 > 500 (FarDistance)
			expectedResult: "",
		},
		{
			name:           "Nearby (Changed)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{100, 100},
			receiverPos:    [2]int32{50, 0}, // Distance 50 <= 100 (NearbyDistance)
			changed:        true,
			expectedResult: "user2",
		},
		{
			name:           "Nearby (Unchanged, Recently Sent)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{50, 0},
			changed:        false,
			lastSent:       200 * time.Millisecond, // < FarInterval (500ms)
			expectedResult: "",
		},
		{
			name:           "Nearby (Unchanged, Sent Long Ago)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{50, 0},
			changed:        false,
			lastSent:       600 * time.Millisecond, // > FarInterval (500ms)
			expectedResult: "user2",
		},
		{
			name:           "Medium Distance (Changed, Below Dynamic Interval)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{300, 0}, // Distance 300. Interval should be 300ms (interpolated between 100 and 500)
			changed:        true,
			lastSent:       100 * time.Millisecond, // < 300ms
			expectedResult: "",
		},
		{
			name:           "Medium Distance (Changed, Above Dynamic Interval)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{300, 0}, // Distance 300. Interval should be 300ms
			changed:        true,
			lastSent:       350 * time.Millisecond, // > 300ms
			expectedResult: "user2",
		},
		{
			name:           "Medium Distance (Unchanged, Below Far Interval)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{300, 0},
			changed:        false,
			lastSent:       400 * time.Millisecond, // < 500ms
			expectedResult: "",
		},
		{
			name:           "Medium Distance (Unchanged, Above Far Interval)",
			senderID:       "user1",
			receiverID:     "user2",
			senderPos:      [2]int32{0, 0},
			receiverPos:    [2]int32{300, 0},
			changed:        false,
			lastSent:       600 * time.Millisecond, // > 500ms
			expectedResult: "user2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			senderEntity := NewUserEntity(tt.senderPos[0], tt.senderPos[1], nil)
			senderEntity.changedAfterSend = tt.changed
			if tt.lastSent > 0 {
				senderEntity.m[tt.receiverID] = time.Now().Add(-tt.lastSent)
			}

			receiverEntity := NewUserEntity(tt.receiverPos[0], tt.receiverPos[1], nil)

			got := m.processSenderReceiverPair(tt.senderID, senderEntity, tt.receiverID, receiverEntity)
			if got != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, got)
			}
		})
	}
}
