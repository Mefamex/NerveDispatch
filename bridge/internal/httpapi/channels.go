package httpapi

import (
	"net/http"

	"whatsapp-bridge/internal/channels"
)

// SetChannelsManager attaches the omnichannel adapter manager to the server.
func (s *Server) SetChannelsManager(mgr *channels.Manager) {
	s.channels = mgr
}

func (s *Server) channelsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var statuses []channels.ChannelStatus

	// Her zaman WhatsApp durumunu ekle
	waStatus := channels.ChannelStatus{
		Name:      channels.ChannelWhatsApp,
		Connected: s.wa.Connected(),
		Paired:    s.wa.IsLoggedIn(),
	}

	if s.channels != nil {
		statuses = s.channels.ListStatus(ctx)
		// Eğer whatsapp listesinde yoksa ekle
		hasWA := false
		for _, st := range statuses {
			if st.Name == channels.ChannelWhatsApp {
				hasWA = true
				break
			}
		}
		if !hasWA {
			statuses = append([]channels.ChannelStatus{waStatus}, statuses...)
		}
	} else {
		statuses = []channels.ChannelStatus{
			waStatus,
			{Name: channels.ChannelTelegram, Connected: false, Paired: false},
			{Name: channels.ChannelDiscord, Connected: false, Paired: false},
			{Name: channels.ChannelInstagram, Connected: false, Paired: false},
		}
	}

	writeJSON(w, map[string]any{
		"status":   "ok",
		"channels": statuses,
	})
}
