package tui

import (
	"context"
	"fmt"
	"github.com/f1bonacc1/process-compose/src/config"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

func (pv *pcView) createStatTable() *tview.Table {
	table := tview.NewTable().SetBorders(false).SetSelectable(false, false)

	table.SetCell(0, 0, tview.NewTableCell("Version:").
		SetSelectable(false))
	table.SetCell(0, 1, tview.NewTableCell(config.Version).
		SetSelectable(false).SetExpansion(1))

	table.SetCell(1, 0, tview.NewTableCell(pv.getHostNameTitle()).
		SetSelectable(false))
	hostname := pv.getHostName()
	table.SetCell(1, 1, tview.NewTableCell(hostname).
		SetSelectable(false).
		SetExpansion(1))

	table.SetCell(2, 0, tview.NewTableCell("Processes:").
		SetSelectable(false))
	pv.procCountCell = tview.NewTableCell(strconv.Itoa(len(pv.procNames))).
		SetSelectable(false).
		SetExpansion(1)
	table.SetCell(2, 1, pv.procCountCell)
	table.SetCell(0, 2, tview.NewTableCell("").
		SetSelectable(false).
		SetExpansion(0))
	table.SetCell(1, 2, tview.NewTableCell("").
		SetSelectable(false).
		SetExpansion(0))

	table.SetCell(0, 3, tview.NewTableCell(pv.getPcTitle()).
		SetSelectable(false).
		SetAlign(tview.AlignRight).
		SetExpansion(1))
	return table
}

func (pv *pcView) getPcTitle() string {
	if pv.project.IsRemote() {
		return config.RemoteProjectName
	} else {
		return config.ProjectName
	}
}

func (pv *pcView) getHostName() string {
	name, err := pv.project.GetHostName()
	if err != nil {
		log.Err(err).Msg("Unable to retrieve hostname")
		return "Unknown"
	}
	return name
}

func (pv *pcView) getHostNameTitle() string {
	if pv.project.IsRemote() {
		return "Server Name:"
	} else {
		return "Hostname:"
	}
}

func (pv *pcView) attentionMessage(message string, duration time.Duration) {
	if duration == 0 {
		return
	}
	go func() {
		pv.appView.QueueUpdateDraw(func() {
			pv.statTable.SetCell(0, 2, tview.NewTableCell(message).
				SetSelectable(false).
				SetAlign(tview.AlignCenter).
				SetExpansion(0).
				SetTextColor(tview.Styles.ContrastSecondaryTextColor).
				SetBackgroundColor(tview.Styles.MoreContrastBackgroundColor))
		})
		time.Sleep(duration)
		pv.hideAttentionMessage()
	}()
}

func (pv *pcView) showAutoProgress(ctx context.Context, duration time.Duration) {
	if duration == 0 {
		return
	}

	full := 10
	step := 1
	go func() {
		ticker := time.NewTicker(duration / time.Duration(full))
		defer ticker.Stop()
		defer pv.statTable.SetCell(1, 2, tview.NewTableCell(""))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				progStr := ""
				if step > full {
					progStr = fmt.Sprintf("%s%s",
						strings.Repeat("□", step-full),
						strings.Repeat("■", 2*full-step))
				} else {
					progStr = fmt.Sprintf("%s%s",
						strings.Repeat("■", step),
						strings.Repeat("□", full-step))
				}
				pv.appView.QueueUpdateDraw(func() {
					pv.statTable.GetCell(1, 2).SetText(progStr)
				})

				step += 1
				if step > 2*full {
					step = 1
				}
			}
		}
	}()
}

func (pv *pcView) hideAttentionMessage() {
	pv.statTable.SetCell(0, 2, tview.NewTableCell(""))
}
