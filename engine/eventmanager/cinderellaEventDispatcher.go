package eventmanager

import "github.com/go-cinderella/cinderella-engine/engine/internal/errs"

type CinderellaEventDispatcher struct {
	EventListeners []CinderellaEventListener
}

func (eventDispatcher *CinderellaEventDispatcher) AddEventListener(listenerToAdd CinderellaEventListener) (err error) {
	if listenerToAdd == nil {
		err = errs.CinderellaError{Msg: "Listener cannot be null."}
	}
	eventDispatcher.EventListeners = append(eventDispatcher.EventListeners, listenerToAdd)
	return err
}

func (eventDispatcher *CinderellaEventDispatcher) DispatchEvent(event CinderellaEvent) (err error) {
	// Call global listeners
	for _, listener := range eventDispatcher.EventListeners {
		if err = listener.OnEvent(event); err != nil {
			return err
		}
	}
	return nil
}
