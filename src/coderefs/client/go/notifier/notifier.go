package notifier

/*
ParseAndNotify parses notifier packets in part of echo/heartbeat and triggers passed callbacks
returns rest of the bytes by removing notifier bytes
*/
func ParseAndNotify(payload []byte, callbacks []func(uint8, uint16, []byte)) []byte {
	return nil
}
