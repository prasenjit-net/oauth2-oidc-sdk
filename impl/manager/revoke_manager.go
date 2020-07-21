package manager

import "net/http"

func (d *DefaultManager) ProcessRevocationEP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	if equestContext, iError := d.RevocationRequestContextFactory(request); iError != nil {
		err := d.JsonErrorWriter(iError, writer, request)
		if err != nil {
			d.ErrorStrategy(err, writer)
		}
		return
	} else {
		for _, handler := range d.RevocationEPHandlers {
			if iError := handler.HandleRevocationEP(ctx, equestContext); iError != nil {
				err := d.JsonErrorWriter(iError, writer, request)
				if err != nil {
					d.ErrorStrategy(err, writer)
				}
				return
			}
		}
		if err := d.RevocationResponseWriter(equestContext, writer, request); err != nil {
			d.ErrorStrategy(err, writer)
		}
	}
}