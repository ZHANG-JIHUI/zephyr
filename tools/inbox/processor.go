package inbox

type Processor interface {
	Invoke(envelopes []Envelope)
}
