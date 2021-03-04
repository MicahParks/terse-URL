package configure

import (
	"io"
	"io/ioutil"

	"github.com/go-openapi/runtime"
	"go.uber.org/zap"
)

// HTMLProducer creates a runtime.ProducerFunc that will read an HTML io.ReadCloser and write it to the given writer.
func HTMLProducer(logger *zap.SugaredLogger) (htmlProducer runtime.ProducerFunc) {
	return func(writer io.Writer, inter interface{}) (err error) {

		// Assume the interface is an io.ReadCloser.
		readCloser, ok := inter.(io.ReadCloser)
		if !ok {
			logger.Warn("Expected HTML data was not io.ReadCloser.")
			return nil
		}

		// Read all the HTML data to write.
		var writeThis []byte
		writeThis, err = ioutil.ReadAll(readCloser)
		if err != nil {
			return err
		}
		if err = readCloser.Close(); err != nil {
			return err
		}

		// Write all the HTML data to the writer.
		if _, err = writer.Write(writeThis); err != nil {
			return err
		}

		return nil
	}
}
