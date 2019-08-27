package cliedit

import (
	"github.com/elves/elvish/cli"
	"github.com/elves/elvish/cli/addons/histlist"
	"github.com/elves/elvish/cli/addons/lastcmd"
	"github.com/elves/elvish/cli/addons/location"
	"github.com/elves/elvish/cli/histutil"
	"github.com/elves/elvish/eval"
	"github.com/elves/elvish/store/storedefs"
)

//elvdoc:fn listing:close
//
// Closes the listing.

func closeListing(app *cli.App) {
	app.MutateAppState(func(s *cli.State) { s.Listing = nil })
}

func initListings(app *cli.App, ev *eval.Evaler, ns eval.Ns, st storedefs.Store) {
	var histStore histutil.Store
	histFuser, err := histutil.NewFuser(st)
	if err == nil {
		histStore = fuserWrapper{histFuser}
	}
	dirStore := dirStore{ev}

	// Common binding and the listing: module.
	lsMap := newBindingVar(emptyBindingMap)
	ns.AddNs("listing",
		eval.Ns{
			"binding": lsMap,
		}.AddGoFns("<edit:listing>:", map[string]interface{}{
			"close": func() { closeListing(app) },
			/*
				"up":               cli.ListingUp,
				"down":             cli.ListingDown,
				"up-cycle":         cli.ListingUpCycle,
				"down-cycle":       cli.ListingDownCycle,
				"toggle-filtering": cli.ListingToggleFiltering,
				"accept":           cli.ListingAccept,
				"accept-close":     cli.ListingAcceptClose,
				"default":          cli.ListingDefault,
			*/
		}))

	histlistMap := newBindingVar(emptyBindingMap)
	histlistBinding := newMapBinding(app, ev, histlistMap, lsMap)
	ns.AddNs("histlist",
		eval.Ns{
			"binding": histlistMap,
		}.AddGoFn("<edit:histlist>", "start", func() {
			histlist.Start(app, histlist.Config{histlistBinding, histStore})
		}))

	lastcmdMap := newBindingVar(emptyBindingMap)
	lastcmdBinding := newMapBinding(app, ev, lastcmdMap, lsMap)
	ns.AddNs("lastcmd",
		eval.Ns{
			"binding": lastcmdMap,
		}.AddGoFn("<edit:lastcmd>", "start", func() {
			// TODO: Specify wordifier
			lastcmd.Start(app, lastcmd.Config{lastcmdBinding, histStore, nil})
		}))

	locationMap := newBindingVar(emptyBindingMap)
	locationBinding := newMapBinding(app, ev, locationMap, lsMap)
	ns.AddNs("location",
		eval.Ns{
			"binding": locationMap,
		}.AddGoFn("<edit:location>", "start", func() {
			location.Start(app, location.Config{locationBinding, dirStore})
		}))
}

// Wraps the histutil.Fuser interface to implement histutil.Store. This is a
// bandaid as we cannot change the implementation of Fuser without breaking its
// other users. Eventually Fuser should implement Store directly.
type fuserWrapper struct {
	*histutil.Fuser
}

func (f fuserWrapper) AddCmd(cmd histutil.Entry) (int, error) {
	return f.Fuser.AddCmd(cmd.Text)
}

// Wraps an Evaler to implement the cli.DirStore interface.
type dirStore struct {
	ev *eval.Evaler
}

func (d dirStore) Chdir(path string) error {
	return d.ev.Chdir(path)
}

func (d dirStore) Dirs() ([]storedefs.Dir, error) {
	return d.ev.DaemonClient.Dirs(map[string]struct{}{})
}
