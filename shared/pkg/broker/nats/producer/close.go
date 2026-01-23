package nats

func (b *NatsBroker) Close() {
	if b == nil || b.conn == nil {
		return
	}
	b.conn.Close()
}
