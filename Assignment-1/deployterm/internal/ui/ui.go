package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aryanmaurya1/deployterm/internal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"gopkg.in/yaml.v3"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	NAME_NAMESPACE_PAGE = "NAMESPACE_PAGE"
)

type stack []string

func (stk *stack) push(element string) {
	*stk = append(*stk, element)
}

func (stk *stack) pop() string {
	if len(*stk) == 0 {
		return ""
	}

	v := (*stk)[len(*stk)-1]
	*stk = (*stk)[:len(*stk)-1]
	return v
}

var stk stack

func switchToErrorPage(app *tview.Application, pages *tview.Pages, err error) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())
	errorModal := tview.NewModal().
		SetText(fmt.Sprintf("Error\n%+v", err.Error())).
		AddButtons([]string{"Back", "Quit"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)

				pages.SwitchToPage(stk.pop())
			} else {
				app.Stop()
			}
		}).
		SetBackgroundColor(tcell.ColorGreen)
	pages.AddAndSwitchToPage(currentPageName, errorModal, false)
}

func switchToEditPage(app *tview.Application, pages *tview.Pages, namespace string, deploymentName string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	deployment, err := opsClient.GetDeployment(context.Background(), namespace, deploymentName)
	if err != nil {
		stk.push(currentPageName)
		switchToErrorPage(app, pages, err)
		return
	}
	deployment.ManagedFields = nil

	jsonData, err := yaml.Marshal(&deployment)
	if err != nil {
		stk.push(currentPageName)
		switchToErrorPage(app, pages, err)
		return
	}

	textArea := tview.NewTextArea().SetWrap(false).SetPlaceholder("Enter text here...")
	textArea.SetTitle(fmt.Sprintf(" Namespace - <%s> | Deployment - <%s> | Edit [pink](press 'Esc' to go back) ", namespace, deploymentName)).SetBorder(true)
	textArea.SetTitleColor(tcell.ColorAntiqueWhite)
	textArea.SetText(string(jsonData), false)

	helpInfo := tview.NewTextView().SetText(" press <ctrl + s> to apply")
	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfos)
	updateInfos()

	mainView := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textArea, 0, 0, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 1, 1, 1, 1, 0, 0, false)

	mainView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Remove dynamically added page
			pages.RemovePage(currentPageName)
			pages.SwitchToPage(stk.pop())
		} else if event.Key() == tcell.KeyCtrlS {
			text := textArea.GetText()
			err = yaml.Unmarshal([]byte(text), &deployment)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return event
			}

			_, err = opsClient.UpdateDeployment(context.Background(), namespace, deployment)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return event
			}
			pages.SwitchToPage(stk.pop())
		}
		return event
	})

	pages.AddAndSwitchToPage(currentPageName, mainView, true)
}

func switchToDetailsPage(app *tview.Application, pages *tview.Pages, namespace string, deploymentName string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle(fmt.Sprintf(" Namespace - <%s> | Deployment - <%s> | Details [pink](press 'Esc' to go back) ", namespace, deploymentName))
	textView.SetTitleColor(tcell.ColorAntiqueWhite)
	textView.SetTextColor(tcell.ColorDarkRed)
	textView.SetFocusFunc(
		func() {
			deployment, err := opsClient.GetDeployment(context.Background(), namespace, deploymentName)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return
			}

			deployment.ManagedFields = nil
			jsonData, err := yaml.Marshal(&deployment)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return
			}

			go func() {
				textView.Clear()
				fmt.Fprintf(textView, "%s", jsonData)
			}()
		})

	textView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEsc {
			// Remove dynamically added page
			pages.RemovePage(currentPageName)
			pages.SwitchToPage(stk.pop())
		}
	})
	pages.AddAndSwitchToPage(currentPageName, textView, true)
}

func switchToOptionsPage(app *tview.Application, pages *tview.Pages, namespace string, deploymentName string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf(" Namespace - <%s> | Deployment - <%s> ", namespace, deploymentName))
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()

			list.AddItem("Describe", "Press to view deployment details", rune('a'), func() {
				stk.push(currentPageName)
				switchToDetailsPage(app, pages, namespace, deploymentName, opsClient)
			})

			list.AddItem("Edit", "Press to edit deployment", rune('b'), func() {
				stk.push(currentPageName)
				switchToEditPage(app, pages, namespace, deploymentName, opsClient)
			})

			list.AddItem("[darkmagenta]Delete", "Press to delete deployment", rune('c'), func() {
				_, err := opsClient.DeleteDeployment(context.Background(), namespace, deploymentName)
				if err != nil {
					stk.push(currentPageName)
					switchToErrorPage(app, pages, err)
					return
				}

				// Remove dynamically added page
				pages.RemovePage(currentPageName)
				pages.SwitchToPage(stk.pop())
			})

			list.AddItem("[red]Back", "Press to go back", rune('x'), func() {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)
				pages.SwitchToPage(stk.pop())
			})

			list.AddItem("[red]Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})

	pages.AddAndSwitchToPage(currentPageName, list, true)
}

func switchToCreateDeploymentPage(app *tview.Application, pages *tview.Pages, namespace string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	textArea := tview.NewTextArea().SetWrap(false).SetPlaceholder("Enter text here...")
	textArea.SetTitle(fmt.Sprintf(" Namespace - <%s> | Create Deployment [pink](press 'Esc' to go back) ", namespace)).SetBorder(true)
	textArea.SetTitleColor(tcell.ColorAntiqueWhite)

	helpInfo := tview.NewTextView().SetText(" press <ctrl + s> to apply")
	position := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)
	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := textArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			position.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			position.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}

	textArea.SetMovedFunc(updateInfos)
	updateInfos()

	mainView := tview.NewGrid().
		SetRows(0, 1).
		AddItem(textArea, 0, 0, 1, 2, 0, 0, true).
		AddItem(helpInfo, 1, 0, 1, 1, 0, 0, false).
		AddItem(position, 1, 1, 1, 1, 0, 0, false)

	mainView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			// Remove dynamically added page
			pages.RemovePage(currentPageName)
			pages.SwitchToPage(stk.pop())
		} else if event.Key() == tcell.KeyCtrlS {
			text := textArea.GetText()

			// Converting YAML to JSON, direct YAML unmarshalling into Deployment object
			// gives error.
			var deploymentMap = make(map[string]any)
			err := yaml.Unmarshal([]byte(text), &deploymentMap)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return event
			}

			bytes, err := json.Marshal(deploymentMap)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
			}

			var deployment = appsv1.Deployment{}
			err = json.Unmarshal(bytes, &deployment)
			if err != nil {
				if err != nil {
					stk.push(currentPageName)
					switchToErrorPage(app, pages, err)
					return event
				}
			}

			_, err = opsClient.CreateDeployment(context.Background(), namespace, &deployment)
			if err != nil {
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
				return event
			}
			pages.SwitchToPage(stk.pop())
		}
		return event
	})

	pages.AddAndSwitchToPage(currentPageName, mainView, true)
}

func switchToDeploymentListPage(app *tview.Application, pages *tview.Pages, namespace string, opsClient internal.IK8sOperation) {
	currentPageName := fmt.Sprintf("%d", time.Now().UnixNano())

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(fmt.Sprintf(" Namespace - <%s> ", namespace))
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()
			deployments, err := opsClient.ListDeployments(context.Background(), namespace)
			if err != nil {
				fmt.Printf("error in listing deployments: %+v\n", err)
				stk.push(currentPageName)
				switchToErrorPage(app, pages, err)
			}

			for idx, deployment := range deployments {
				primary := deployment.Name
				secondary := fmt.Sprintf("desired_replicas: %+v | current_replica: %+v", *deployment.Spec.Replicas, deployment.Status.Replicas)
				shortcut := rune(int('a') + idx)

				// Avoid shortcut collision with Quit/Back operations
				if shortcut == 'q' {
					shortcut = '1'
				} else if shortcut == 'x' {
					shortcut = '2'
				} else if shortcut == 'n' {
					shortcut = '3'
				}

				list.AddItem(primary, secondary, shortcut, func() {
					stk.push(currentPageName)
					switchToOptionsPage(app, pages, namespace, deployment.Name, opsClient)
				})
			}

			list.AddItem("[pink]Create New", "Press to create new deployment", rune('n'), func() {
				stk.push(currentPageName)
				switchToCreateDeploymentPage(app, pages, namespace, opsClient)
			})

			list.AddItem("[red]Back", "Press to go back", rune('x'), func() {
				// Remove dynamically added page
				pages.RemovePage(currentPageName)
				pages.SwitchToPage(stk.pop())
			})

			list.AddItem("[red]Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})

	pages.AddAndSwitchToPage(currentPageName, list, true)
}

func namespaceListPage(app *tview.Application, pages *tview.Pages, opsClient internal.IK8sOperation) tview.Primitive {
	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle(" Namespaces ")
	list.SetTitleColor(tcell.ColorAntiqueWhite)
	list.SetFocusFunc(
		func() {
			list.Clear()
			namespaces, err := opsClient.ListNamespaces(context.Background())
			if err != nil {
				fmt.Printf("error in listing namespace: %+v\n", err)
				stk.push(NAME_NAMESPACE_PAGE)
				switchToErrorPage(app, pages, err)
			}

			for idx, ns := range namespaces {
				primary := ns.Name
				secondary := fmt.Sprintf("create_at: %+v", ns.CreationTimestamp)
				shortcut := rune(int('a') + idx)

				// Avoid shortcut collision with Quit/Back operations
				if shortcut == 'q' {
					shortcut = '1'
				}

				list.AddItem(primary, secondary, shortcut, func() {
					stk.push(NAME_NAMESPACE_PAGE)
					switchToDeploymentListPage(app, pages, ns.Name, opsClient)
				})
			}

			list.AddItem("[red]Quit", "Press to exit", rune('q'), func() {
				app.Stop()
			})

		})
	return list
}

func GetRootPage(app *tview.Application, pages *tview.Pages, opsClient internal.IK8sOperation) (tview.Primitive, string) {
	return namespaceListPage(app, pages, opsClient), NAME_NAMESPACE_PAGE
}
