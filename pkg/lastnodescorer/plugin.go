package lastnodescorer

import (
	"context"
	"sort"
	"k8s.io/klog/v2"
	v1 "k8s.io/api/core/v1"
	fwk "k8s.io/kube-scheduler/framework"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"k8s.io/apimachinery/pkg/runtime"
)

type LastNodeScorer struct {
	logger klog.Logger
	handle framework.Handle
}

var _ = framework.ScorePlugin(&LastNodeScorer{})

const Name = "LastNodeScorer"

func (l *LastNodeScorer) Name() string {
	return Name
}

func (l *LastNodeScorer) Score(ctx context.Context, state fwk.CycleState, pod *v1.Pod, nodeInfo fwk.NodeInfo) (int64, *fwk.Status) {
	//logger := klog.FromContext(klog.NewContext(ctx, l.logger)).WithValues("ExtensionPoint", "Score")
	//logger.V(10).Info("No match for mode", "mode", mode)
	// 🔥 Forma correta de acessar nós
	nodeInfos, err := l.handle.SnapshotSharedLister().NodeInfos().List()
	if err != nil {
		return 0, fwk.AsStatus(err)
	}

	var names []string
	for _, n := range nodeInfos {
		names = append(names, n.Node().Name)
	}

	sort.Strings(names)
	last := names[len(names)-1]

	current := nodeInfo.Node().Name

	if current == last {
		return framework.MaxNodeScore, fwk.NewStatus(fwk.Success)
	}

	return 0, fwk.NewStatus(fwk.Success)
}

func (l *LastNodeScorer) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

func New(ctx context.Context, _ runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	logger := klog.FromContext(ctx).WithValues("plugin", Name)
	return &LastNodeScorer{
		handle: handle,
		logger: logger,
	}, nil
}