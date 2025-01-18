package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// private channel to force all pipes to close
	pipelineTermination := make(Bi)
	nextIn := in

	for _, stage := range stages {
		stageOut := stage(nextIn)
		nextIn = startPipe(stageOut, pipelineTermination, nil)
	}

	// last pipe closes all pipes if done channel is triggered
	return startPipe(nextIn, done, func() {
		close(pipelineTermination)
	})
}

func startPipe(stageOut Out, closeThisPipe In, onPipeClose func()) Out {
	pipeCh := make(Bi)

	go func() {
		defer func() {
			if onPipeClose != nil {
				onPipeClose()
			}
			close(pipeCh)

			// drain the channel to unblock stage that might be stuck
			// trying to send messages to this channel
			for range stageOut { //revive:disable-line:empty-block
				// ignore remaining messages
			}
		}()

		for {
			select {
			case msg, ok := <-stageOut:
				{
					if !ok {
						// if stage completed, close this pipe
						return
					}

					select {
					// attempt to send the message into next stage via pipeCh
					case pipeCh <- msg:
					case <-closeThisPipe:
						// if the termination channel is triggered, close this pipe instead
						return
					}
				}
			case <-closeThisPipe:
				// if the termination channel is triggered, close this pipe instead
				return
			}
		}
	}()

	return pipeCh
}
