package callbacks

import "github.com/sirupsen/logrus"

// PrintErrorMessage prints an error message using logrus
func PrintErrorMessage(err error) {
    logrus.WithError(err).Error("Operation failed")
}
