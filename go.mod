module VideoControls

go 1.15

require (
	fyne.io/fyne v1.4.0

	// Use branch "forked" of janpfeifer/webcam, until changes make into github.com/blackjack/webcam
	github.com/janpfeifer/webcam v0.0.0-20201102083701-f80efa8fc9e2
)

// When developing with local copy of `webcam`, uncomment the following.
// replace github.com/janpfeifer/webcam => ../webcam
