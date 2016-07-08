package mesosfmt

import (
	"github.com/golang/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSuffix(t *testing.T) {
	sfx := suffix("asdqwe", 0)
	assert.Equal(t, sfx, "")

	sfx = suffix("asdqwe", 3)
	assert.Equal(t, sfx, "qwe")

	sfx = suffix("asdqwe", 6)
	assert.Equal(t, sfx, "asdqwe")

	sfx = suffix("asdqwe", 10)
	assert.Equal(t, sfx, "asdqwe")
}

func TestID(t *testing.T) {
	id := ID("487c73d8-9951-f23c-34bd-8085bfd30c49")
	assert.Equal(t, id, "#30c49")
}

func TestResource(t *testing.T) {
	mem := Resource(util.NewScalarResource("mem", 512))
	assert.Equal(t, mem, "mem:512.00")

	ports := Resource(util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(31000, 32000)}))
	assert.Equal(t, ports, "ports:[31000..32000]")

	ports = Resource(util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(4000, 7000), util.NewValueRange(31000, 32000)}))
	assert.Equal(t, ports, "ports:[4000..7000][31000..32000]")
}

func TestResources(t *testing.T) {
	resources := Resources([]*mesos.Resource{util.NewScalarResource("cpus", 4), util.NewScalarResource("mem", 512), util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(31000, 32000)})})
	assert.Contains(t, resources, "cpus")
	assert.Contains(t, resources, "mem")
	assert.Contains(t, resources, "ports")
}

func TestAttribute(t *testing.T) {
	attr := Attribute(&mesos.Attribute{
		Name:   proto.String("rack"),
		Type:   mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: proto.Float64(2)},
	})
	assert.Equal(t, attr, "rack:2.00")

	attr = Attribute(&mesos.Attribute{
		Name: proto.String("datacenter"),
		Type: mesos.Value_TEXT.Enum(),
		Text: &mesos.Value_Text{Value: proto.String("DC-1")},
	})
	assert.Equal(t, attr, "datacenter:DC-1")
}

func TestAttributes(t *testing.T) {
	attributes := Attributes([]*mesos.Attribute{{
		Name:   proto.String("rack"),
		Type:   mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: proto.Float64(2)},
	}, {
		Name:   proto.String("dc"),
		Type:   mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: proto.Float64(1)},
	}})
	assert.Contains(t, attributes, "rack")
	assert.Contains(t, attributes, "dc")
}

func TestOffer(t *testing.T) {
	offer := util.NewOffer(util.NewOfferID("487c73d8-9951-f23c-34bd-8085bfd30c49"), util.NewFrameworkID("20150903-065451-84125888-5050-10715-0053"),
		util.NewSlaveID("20150903-065451-84125888-5050-10715-S1"), "slave0")
	assert.Equal(t, Offer(offer), "slave0#30c49")

	offer.Resources = []*mesos.Resource{util.NewScalarResource("cpus", 4), util.NewScalarResource("mem", 512), util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(31000, 32000)})}
	assert.Equal(t, Offer(offer), "slave0#30c49 cpus:4.00 mem:512.00 ports:[31000..32000]")

	offer.Attributes = []*mesos.Attribute{{
		Name:   proto.String("rack"),
		Type:   mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: proto.Float64(2)},
	}}
	assert.Equal(t, Offer(offer), "slave0#30c49 cpus:4.00 mem:512.00 ports:[31000..32000] rack:2.00")

	offer.Resources = nil
	assert.Equal(t, Offer(offer), "slave0#30c49 rack:2.00")
}

func TestOffers(t *testing.T) {
	offer1 := util.NewOffer(util.NewOfferID("487c73d8-9951-f23c-34bd-8085bfd30c49"), util.NewFrameworkID("20150903-065451-84125888-5050-10715-0053"),
		util.NewSlaveID("20150903-065451-84125888-5050-10715-S1"), "slave0")
	offer1.Resources = []*mesos.Resource{util.NewScalarResource("cpus", 4), util.NewScalarResource("mem", 512), util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(31000, 32000)})}

	offer2 := util.NewOffer(util.NewOfferID("26d5b34c-ef81-638d-5ad5-32c743c9c033"), util.NewFrameworkID("20150903-065451-84125888-5050-10715-0037"),
		util.NewSlaveID("20150903-065451-84125888-5050-10715-S0"), "master")
	offer2.Resources = []*mesos.Resource{util.NewScalarResource("cpus", 2), util.NewScalarResource("mem", 1024), util.NewRangesResource("ports", []*mesos.Value_Range{util.NewValueRange(4000, 7000)})}
	offer2.Attributes = []*mesos.Attribute{{
		Name:   proto.String("rack"),
		Type:   mesos.Value_SCALAR.Enum(),
		Scalar: &mesos.Value_Scalar{Value: proto.Float64(2)},
	}}

	offers := Offers([]*mesos.Offer{offer1, offer2})
	assert.Len(t, strings.Split(offers, "\n"), 2)
}

func TestStatus(t *testing.T) {
	status := util.NewTaskStatus(util.NewTaskID("task-1"), mesos.TaskState_TASK_FINISHED)
	assert.Contains(t, Status(status), "task-1 TASK_FINISHED reason")

	status.SlaveId = util.NewSlaveID("20150903-065451-84125888-5050-10715-S1")
	assert.Contains(t, Status(status), "task-1 TASK_FINISHED slave: #15-S1")

	status.State = mesos.TaskState_TASK_RUNNING.Enum()
	assert.NotContains(t, Status(status), "reason")

	status.State = mesos.TaskState_TASK_LOST.Enum()
	status.Reason = mesos.TaskStatus_REASON_EXECUTOR_TERMINATED.Enum()
	assert.Contains(t, Status(status), "task-1 TASK_LOST slave: #15-S1 reason: REASON_EXECUTOR_TERMINATED")

	status.Message = proto.String("boom!")
	assert.Contains(t, Status(status), "message: boom!")
}
