// components contain layout code that are shared by multiple pages but aren't widely used enough to be defined as
// widgets

package ui

import (
	"image"

	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/decred/dcrd/dcrutil"
	"github.com/planetdecred/dcrlibwallet"
	"github.com/planetdecred/godcr/ui/decredmaterial"
	"github.com/planetdecred/godcr/ui/values"
	"github.com/planetdecred/godcr/wallet"
)

// layoutBalance aligns the main and sub DCR balances horizontally, putting the sub
// balance at the baseline of the row.
func (page pageCommon) layoutBalance(gtx layout.Context, amount string) layout.Dimensions {
	// todo: make "DCR" symbols small when there are no decimals in the balance
	mainText, subText := breakBalance(page.printer, amount)
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Baseline}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(values.TextSize20, mainText).Layout(gtx)
		}),
		layout.Rigid(func(gtx C) D {
			return page.theme.Label(values.TextSize14, subText).Layout(gtx)
		}),
	)
}

// layoutTopBar is the top horizontal bar on every page of the app. It lays out the wallet balance, receive and send
// buttons.
func (page pageCommon) layoutTopBar(gtx layout.Context) layout.Dimensions {
	card := page.theme.Card()
	card.Radius = decredmaterial.CornerRadius{}
	return card.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						return layout.W.Layout(gtx, func(gtx C) D {
							h := values.MarginPadding24
							v := values.MarginPadding16
							// Balance container
							return Container{padding: layout.Inset{Right: h, Left: h, Top: v, Bottom: v}}.Layout(gtx,
								func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											img := page.icons.logo
											img.Scale = 1.0
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx,
												func(gtx C) D {
													return img.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Center.Layout(gtx, func(gtx C) D {
												return page.layoutBalance(gtx, page.info.TotalBalance)
											})
										}),
									)
								})
						})
					}),
					layout.Rigid(func(gtx C) D {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Inset{Right: values.MarginPadding8}.Layout(gtx, func(gtx C) D {
								list := layout.List{Axis: layout.Horizontal}
								return list.Layout(gtx, len(page.appBarNavItems), func(gtx C, i int) D {
									// header buttons container
									return Container{layout.UniformInset(values.MarginPadding16)}.Layout(gtx, func(gtx C) D {
										return decredmaterial.Clickable(gtx, page.appBarNavItems[i].clickable, func(gtx C) D {
											return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
												layout.Rigid(func(gtx C) D {
													return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
														func(gtx C) D {
															return layout.Center.Layout(gtx, func(gtx C) D {
																img := page.appBarNavItems[i].image
																img.Scale = 1.0
																return page.appBarNavItems[i].image.Layout(gtx)
															})
														})
												}),
												layout.Rigid(func(gtx C) D {
													return layout.Inset{
														Left: values.MarginPadding0,
													}.Layout(gtx, func(gtx C) D {
														return layout.Center.Layout(gtx, func(gtx C) D {
															return page.theme.Body1(page.appBarNavItems[i].page).Layout(gtx)
														})
													})
												}),
											)
										})
									})
								})
							})
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Constraints.Max.X
				return page.theme.Separator().Layout(gtx)
			}),
		)
	})
}

const (
	navDrawerWidth          = 160
	navDrawerMinimizedWidth = 72
)

// layoutNavDrawer is the left vertical pane on every page of the app. It vertically lays out buttons used to navigate
// to different pages.
func (page pageCommon) layoutNavDrawer(gtx layout.Context) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx C) D {
			list := layout.List{Axis: layout.Vertical}
			return list.Layout(gtx, len(page.drawerNavItems), func(gtx C, i int) D {
				background := page.theme.Color.Surface
				if page.drawerNavItems[i].page == *page.page {
					background = page.theme.Color.LightGray
				}
				txt := page.theme.Label(values.TextSize16, page.drawerNavItems[i].page)
				return decredmaterial.Clickable(gtx, page.drawerNavItems[i].clickable, func(gtx C) D {
					card := page.theme.Card()
					card.Color = background
					card.Radius = decredmaterial.CornerRadius{NE: 0, NW: 0, SE: 0, SW: 0}
					return card.Layout(gtx, func(gtx C) D {
						return Container{
							layout.Inset{
								Top:    values.MarginPadding16,
								Right:  values.MarginPadding24,
								Bottom: values.MarginPadding16,
								Left:   values.MarginPadding24,
							},
						}.Layout(gtx, func(gtx C) D {
							axis := layout.Horizontal
							leftInset := values.MarginPadding15
							width := navDrawerWidth
							if *page.isNavDrawerMinimized {
								axis = layout.Vertical
								txt.TextSize = values.TextSize10
								leftInset = values.MarginPadding0
								width = navDrawerMinimizedWidth
							}

							gtx.Constraints.Min.X = int(gtx.Metric.PxPerDp) * width
							return layout.Flex{Axis: axis}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									img := page.drawerNavItems[i].imageInactive
									if page.drawerNavItems[i].page == *page.page {
										img = page.drawerNavItems[i].image
									}
									return layout.Center.Layout(gtx, func(gtx C) D {
										img.Scale = 1.0
										return img.Layout(gtx)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Inset{
										Left: leftInset,
										Top:  values.MarginPadding4,
									}.Layout(gtx, func(gtx C) D {
										textColor := page.theme.Color.Gray3
										if page.drawerNavItems[i].page == *page.page {
											textColor = page.theme.Color.DeepBlue
										}
										txt.Color = textColor
										return layout.Center.Layout(gtx, txt.Layout)
									})
								}),
							)
						})
					})
				})
			})
		}),
		layout.Expanded(func(gtx C) D {
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			return layout.SE.Layout(gtx, func(gtx C) D {
				btn := page.minimizeNavDrawerButton
				if *page.isNavDrawerMinimized {
					btn = page.maximizeNavDrawerButton
				}
				return btn.Layout(gtx)
			})
		}),
	)
}

type TransactionRow struct {
	transaction wallet.Transaction
	index       int
	showBadge   bool
}

// transactionRow is a single transaction row on the transactions and overview page. It lays out a transaction's
// direction, balance, status.
func transactionRow(gtx layout.Context, common pageCommon, row TransactionRow) layout.Dimensions {
	gtx.Constraints.Min.X = gtx.Constraints.Max.X
	directionIconTopMargin := values.MarginPadding16

	if row.index == 0 && row.showBadge {
		directionIconTopMargin = values.MarginPadding14
	} else if row.index == 0 {
		// todo: remove top margin from container
		directionIconTopMargin = values.MarginPadding0
	}

	return layout.Inset{Top: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				icon := common.icons.receiveIcon
				if row.transaction.Txn.Direction == dcrlibwallet.TxDirectionSent {
					icon = common.icons.sendIcon
				}
				icon.Scale = 1.0

				return layout.Inset{Top: directionIconTopMargin}.Layout(gtx, func(gtx C) D {
					return icon.Layout(gtx)
				})
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if row.index == 0 {
							return layout.Dimensions{}
						}
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						separator := common.theme.Separator()
						separator.Width = gtx.Constraints.Max.X - gtx.Px(unit.Dp(16))
						return layout.E.Layout(gtx, func(gtx C) D {
							// Todo: add comment
							marginBottom := values.MarginPadding16
							if row.showBadge {
								marginBottom = values.MarginPadding5
							}
							return layout.Inset{Bottom: marginBottom}.Layout(gtx,
								func(gtx C) D {
									return separator.Layout(gtx)
								})
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						gtx.Constraints.Min.X = gtx.Constraints.Max.X
						return layout.Inset{}.Layout(gtx, func(gtx C) D {
							return layout.Flex{
								Axis:      layout.Horizontal,
								Spacing:   layout.SpaceBetween,
								Alignment: layout.Middle,
							}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									return layout.Inset{Left: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
										return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
											layout.Rigid(func(gtx C) D {
												return common.layoutBalance(gtx, row.transaction.Balance)
											}),
											layout.Rigid(func(gtx C) D {
												if row.showBadge {
													return walletLabel(gtx, common, row.transaction.WalletName)
												}
												return layout.Dimensions{}
											}),
										)
									})
								}),
								layout.Rigid(func(gtx C) D {
									return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding8}.Layout(gtx,
												func(gtx C) D {
													s := formatDateOrTime(row.transaction.Txn.Timestamp)
													if row.transaction.Status != "confirmed" {
														s = row.transaction.Status
													}
													status := common.theme.Body1(s)
													status.Color = common.theme.Color.Gray
													status.Alignment = text.Middle
													return status.Layout(gtx)
												})
										}),
										layout.Rigid(func(gtx C) D {
											return layout.Inset{Right: values.MarginPadding16}.Layout(gtx, func(gtx C) D {
												statusIcon := common.icons.confirmIcon
												if row.transaction.Status != "confirmed" {
													statusIcon = common.icons.pendingIcon
												}
												statusIcon.Scale = 1.0
												return statusIcon.Layout(gtx)
											})
										}),
									)
								}),
							)
						})
					}),
				)
			}),
		)
	})
}

// walletLabel displays the wallet which a transaction belongs to. It is only displayed on the overview page when there
// are transactions from multiple wallets
func walletLabel(gtx layout.Context, c pageCommon, walletName string) D {
	return decredmaterial.Card{
		Color: c.theme.Color.LightGray,
	}.Layout(gtx, func(gtx C) D {
		return Container{
			layout.Inset{
				Left:  values.MarginPadding4,
				Right: values.MarginPadding4,
			}}.Layout(gtx, func(gtx C) D {
			name := c.theme.Label(values.TextSize12, walletName)
			name.Color = c.theme.Color.Gray
			return name.Layout(gtx)
		})
	})
}

// endToEndRow layouts out its content on both ends of its horizontal layout.
func endToEndRow(gtx layout.Context, leftWidget, rightWidget func(C) D) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return leftWidget(gtx)
		}),
		layout.Flexed(1, func(gtx C) D {
			return layout.E.Layout(gtx, func(gtx C) D {
				return rightWidget(gtx)
			})
		}),
	)
}

func (page *pageCommon) accountSelectorLayout(gtx layout.Context, title string) layout.Dimensions {
	border := widget.Border{
		Color:        page.theme.Color.Gray1,
		CornerRadius: values.MarginPadding8,
		Width:        values.MarginPadding2,
	}
	page.wallAcctSelector.title = title

	return border.Layout(gtx, func(gtx C) D {
		return layout.UniformInset(values.MarginPadding12).Layout(gtx, func(gtx C) D {
			return decredmaterial.Clickable(gtx, page.wallAcctSelector.fromAccount, func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
					layout.Rigid(func(gtx C) D {
						accountIcon := page.icons.accountIcon
						accountIcon.Scale = 1
						inset := layout.Inset{
							Right: values.MarginPadding8,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Rigid(func(gtx C) D {
						return page.theme.Body1(
							page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount].Name).Layout(gtx)
					}),
					layout.Rigid(func(gtx C) D {
						inset := layout.Inset{
							Left: values.MarginPadding4,
							Top:  values.MarginPadding2,
						}
						return inset.Layout(gtx, func(gtx C) D {
							return decredmaterial.Card{
								Color: page.theme.Color.LightGray,
							}.Layout(gtx, func(gtx C) D {
								m2 := values.MarginPadding2
								m4 := values.MarginPadding4
								inset := layout.Inset{
									Left:   m4,
									Top:    m2,
									Bottom: m2,
									Right:  m4,
								}
								return inset.Layout(gtx, func(gtx C) D {
									text := page.theme.Caption(page.info.Wallets[*page.selectedWallet].Name)
									text.Color = page.theme.Color.Gray
									return text.Layout(gtx)
								})
							})
						})
					}),
					layout.Flexed(1, func(gtx C) D {
						return layout.E.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Rigid(func(gtx C) D {
									txt := page.theme.Body1(
										page.info.Wallets[*page.selectedWallet].Accounts[*page.selectedAccount].TotalBalance)
									txt.Color = page.theme.Color.DeepBlue
									return txt.Layout(gtx)
								}),
								layout.Rigid(func(gtx C) D {
									inset := layout.Inset{
										Left: values.MarginPadding15,
									}
									return inset.Layout(gtx, func(gtx C) D {
										return page.icons.dropDownIcon.Layout(gtx, values.MarginPadding20)
									})
								}),
							)
						})
					}),
				)
			})
		})
	})
}

func (page *pageCommon) walletAccountModalLayout(gtx layout.Context) layout.Dimensions {
	wallAcctGroup := func(gtx layout.Context, title string, body layout.Widget) layout.Dimensions {
		return layout.Inset{
			Bottom: values.MarginPadding10,
		}.Layout(gtx, func(gtx C) D {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					txt := page.theme.Body2(title)
					txt.Color = page.theme.Color.Text
					inset := layout.Inset{
						Bottom: values.MarginPadding15,
					}
					return inset.Layout(gtx, txt.Layout)
				}),
				layout.Rigid(body),
			)
		})
	}
	wallAcctSelector := page.wallAcctSelector
	w := []func(gtx C) D{
		func(gtx C) D {
			tit := page.theme.Label(values.TextSize24, wallAcctSelector.title)
			tit.Color = page.theme.Color.Text
			return tit.Layout(gtx)
		},
		func(gtx C) D {
			return wallAcctSelector.walletsList.Layout(gtx, len(page.info.Wallets), func(gtx C, windex int) D {
				if page.info.Wallets[windex].IsWatchingOnly {
					return D{}
				}

				return wallAcctGroup(gtx, page.info.Wallets[windex].Name, func(gtx C) D {
					return wallAcctSelector.accountsList.Layout(gtx, len(page.info.Wallets[windex].Accounts), func(gtx C, aindex int) D {
						click := wallAcctSelector.walletAccounts[windex][aindex].evt
						pointer.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Add(gtx.Ops)
						click.Add(gtx.Ops)
						page.walletAccountsHandler(gtx, wallAcctSelector.walletAccounts[windex][aindex])
						return page.walletAccountLayout(gtx, wallAcctSelector.walletAccounts[windex][aindex])
					})
				})
			})
		},
	}

	return wallAcctSelector.walletAccount.Layout(gtx, w, 850)
}

func (page *pageCommon) walletAccountsHandler(gtx layout.Context, wallAcct walletAccount) {
	for _, e := range wallAcct.evt.Events(gtx) {
		if e.Type == gesture.TypeClick {
			*page.selectedWallet = wallAcct.walletIndex
			*page.selectedAccount = wallAcct.accountIndex
			page.wallAcctSelector.isWalletAccountModalOpen = false
		}
	}
}

func (page *pageCommon) walletAccountLayout(gtx layout.Context, wallAcct walletAccount) layout.Dimensions {
	accountIcon := page.icons.accountIcon
	accountIcon.Scale = 1
	return layout.Inset{
		Bottom: values.MarginPadding20,
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(0.1, func(gtx C) D {
						return layout.Inset{
							Right: values.MarginPadding18,
						}.Layout(gtx, func(gtx C) D {
							return accountIcon.Layout(gtx)
						})
					}),
					layout.Flexed(0.8, func(gtx C) D {
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
							layout.Rigid(func(gtx C) D {
								acct := page.theme.Label(values.TextSize18, wallAcct.accountName)
								acct.Color = page.theme.Color.Text
								return endToEndRow(gtx, acct.Layout, func(gtx C) D {
									return page.layoutBalance(gtx, wallAcct.totalBalance)
								})
							}),
							layout.Rigid(func(gtx C) D {
								spendable := page.theme.Label(values.TextSize14, "Spendable")
								spendable.Color = page.theme.Color.Gray
								spendableBal := page.theme.Label(values.TextSize14, wallAcct.spendable)
								spendableBal.Color = page.theme.Color.Gray
								return endToEndRow(gtx, spendable.Layout, spendableBal.Layout)
							}),
						)
					}),
					layout.Flexed(0.1, func(gtx C) D {
						if *page.selectedWallet == wallAcct.walletIndex && *page.selectedAccount == wallAcct.accountIndex {
							return layout.E.Layout(gtx, func(gtx C) D {
								return page.icons.navigationCheck.Layout(gtx, values.MarginPadding20)
							})
						}
						return layout.Dimensions{}
					}),
				)
			}),
		)
	})
}

func (page pageCommon) handleNavEvents() {
	for page.minimizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = true
	}

	for page.maximizeNavDrawerButton.Button.Clicked() {
		*page.isNavDrawerMinimized = false
	}

	for i := range page.appBarNavItems {
		for page.appBarNavItems[i].clickable.Clicked() {
			page.changePage(page.appBarNavItems[i].page)
		}
	}

	for i := range page.drawerNavItems {
		for page.drawerNavItems[i].clickable.Clicked() {
			page.changePage(page.drawerNavItems[i].page)
		}
	}

	for windex := 0; windex < page.info.LoadedWallets; windex++ {
		if _, ok := page.wallAcctSelector.walletAccounts[windex]; !ok {
			accounts := page.info.Wallets[windex].Accounts
			if len(accounts) != len(page.wallAcctSelector.walletAccounts[windex]) {
				page.wallAcctSelector.walletAccounts[windex] = make([]walletAccount, len(accounts))
				for aindex := range accounts {
					page.wallAcctSelector.walletAccounts[windex][aindex] = walletAccount{
						walletIndex:  windex,
						accountIndex: aindex,
						evt:          &gesture.Click{},
						accountName:  accounts[aindex].Name,
						totalBalance: accounts[aindex].TotalBalance,
						spendable:    dcrutil.Amount(accounts[aindex].SpendableBalance).String(),
					}
				}
			}
		}
	}

	if page.wallAcctSelector.fromAccount.Clicked() {
		page.wallAcctSelector.isWalletAccountModalOpen = true
	}
}
