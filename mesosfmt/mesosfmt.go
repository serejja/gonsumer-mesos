// Package mesosfmt provides human-readable representations for general Mesos entities.
package mesosfmt

import (
	"bytes"
	"fmt"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"strings"
)

func suffix(str string, maxLen int) string {
	if len(str) < maxLen {
		return str
	}

	return str[len(str)-maxLen:]
}

// ID trims a given string to show only last 5 characters and starts with a # sign.
func ID(id string) string {
	return fmt.Sprintf("#%s", suffix(id, 5))
}

// Offers formats a Mesos offer slice to a human readable format.
func Offers(offers []*mesos.Offer) string {
	var offerStrings = make([]string, len(offers))

	for idx, offer := range offers {
		offerStrings[idx] = Offer(offer)
	}

	return strings.Join(offerStrings, "\n")
}

// Offer formats a single Mesos offer to a human readable format.
func Offer(offer *mesos.Offer) string {
	var buffer bytes.Buffer

	_, _ = buffer.WriteString(offer.GetHostname())
	_, _ = buffer.WriteString(ID(offer.GetId().GetValue()))
	resources := Resources(offer.GetResources())
	if resources != "" {
		_, _ = buffer.WriteString(" ")
		_, _ = buffer.WriteString(resources)
	}
	attributes := Attributes(offer.GetAttributes())
	if attributes != "" {
		_, _ = buffer.WriteString(" ")
		_, _ = buffer.WriteString(attributes)
	}

	return buffer.String()
}

// Resources formats a Mesos resource slice to a human readable format.
func Resources(resources []*mesos.Resource) string {
	var buffer bytes.Buffer

	for _, resource := range resources {
		if buffer.Len() != 0 {
			_, _ = buffer.WriteString(" ")
		}
		_, _ = buffer.WriteString(Resource(resource))
	}

	return buffer.String()
}

// Resource formats a single Mesos resource to a human readable format.
func Resource(resource *mesos.Resource) string {
	var buffer bytes.Buffer

	_, _ = buffer.WriteString(resource.GetName())
	_, _ = buffer.WriteString(":")
	if resource.GetScalar() != nil {
		_, _ = buffer.WriteString(fmt.Sprintf("%.2f", resource.GetScalar().GetValue()))
	}
	if resource.GetRanges() != nil {
		for _, r := range resource.GetRanges().GetRange() {
			_, _ = buffer.WriteString(fmt.Sprintf("[%d..%d]", r.GetBegin(), r.GetEnd()))
		}
	}

	return buffer.String()
}

// Attributes formats a Mesos attribute slice to a human readable format.
func Attributes(attributes []*mesos.Attribute) string {
	var buffer bytes.Buffer

	for _, attr := range attributes {
		if buffer.Len() != 0 {
			_, _ = buffer.WriteString(";")
		}
		_, _ = buffer.WriteString(Attribute(attr))
	}

	return buffer.String()
}

// Attribute formats a single Mesos attribute to a human readable format.
func Attribute(attribute *mesos.Attribute) string {
	var buffer bytes.Buffer

	_, _ = buffer.WriteString(attribute.GetName())
	_, _ = buffer.WriteString(":")
	if attribute.GetText() != nil {
		_, _ = buffer.WriteString(attribute.GetText().GetValue())
	}
	if attribute.GetScalar() != nil {
		_, _ = buffer.WriteString(fmt.Sprintf("%.2f", attribute.GetScalar().GetValue()))
	}

	return buffer.String()
}

// Status formats a Mesos TaskStatus to a human readable format.
func Status(status *mesos.TaskStatus) string {
	var buffer bytes.Buffer
	_, _ = buffer.WriteString(fmt.Sprintf("%s %s", status.GetTaskId().GetValue(), status.GetState().String()))
	if status.GetSlaveId() != nil && status.GetSlaveId().GetValue() != "" {
		_, _ = buffer.WriteString(" slave: ")
		_, _ = buffer.WriteString(ID(status.GetSlaveId().GetValue()))
	}

	if status.GetState() != mesos.TaskState_TASK_RUNNING {
		_, _ = buffer.WriteString(" reason: ")
		_, _ = buffer.WriteString(status.GetReason().String())
	}

	if status.GetMessage() != "" {
		_, _ = buffer.WriteString(" message: ")
		_, _ = buffer.WriteString(status.GetMessage())
	}

	return buffer.String()
}
