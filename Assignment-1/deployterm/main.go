package main

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/aryanmaurya1/deployterm/internal"
	"github.com/aryanmaurya1/deployterm/internal/infra"
	"github.com/aryanmaurya1/deployterm/internal/ui"
	"github.com/rivo/tview"
)

func main() {
	var kconfigPath string
	var useControllerRuntime bool

	flag.StringVar(&kconfigPath, "kubeconfig", "", "path to kubeconfig file")
	flag.BoolVar(&useControllerRuntime, "use-controller-runtime", false, "use controller-runtime library instead of client-go")
	flag.Parse()

	k8sClient, err := internal.NewK8sClient(kconfigPath)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to k8s cluster: %+v\n", err.Error()))
	}

	var opsClient internal.IK8sOperation

	opsClient = infra.NewClientsetWrapper(k8sClient.GetClientset())
	if useControllerRuntime {
		opsClient = infra.NewRunclientWrapper(k8sClient.GetRunclient())
	}

	fmt.Printf("using client: %+v\n", reflect.TypeOf(opsClient))

	app := tview.NewApplication()
	pages := tview.NewPages()

	rootPage, rootPageName := ui.GetRootPage(app, pages, opsClient)
	pages.AddPage(rootPageName, rootPage, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
