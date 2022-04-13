package main

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var PluginName = "wttrin"
var PluginVersion = ""
var PluginBuilddate = ""

const (
	baseURL   = "https://wttr.in/"
	baseURLv2 = "https://v2.wttr.in/"

	curlUserAgent = "curl/7.54.0"
)

var _httpClient *http.Client

func Start(discord *discordgo.Session) {
	discord.AddHandler(onMessageCreate)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(msg, " ")

	channel, _ := s.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}
	guild, _ := s.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}
	if strings.ToLower(parts[0]) == "!wttr" || strings.Contains(strings.ToLower(parts[0]), "!wttrp") {
		handleWttrQuery(s, m, parts, guild)
	}
}

func handleWttrQuery(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) {
	if len(parts) > 1 {
		query := strings.Split(strings.Join(parts[1:], "%20"), "?")
		switch parts[0] {
		case "!wttr":
			wttr, err := weather(query[0] + "?format=4")
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("wttr.in query failed")
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, string(wttr))
		case "!wttrp":
			var wttr []byte
			var err error
			if len(query) > 1 {
				wttr, err = weather(url.QueryEscape(query[0]) + ".png" + "?" + query[1])
			} else {
				wttr, err = weather(url.QueryEscape(query[0]) + ".png?0")
			}
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("wttr.in query failed")
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			s.ChannelFileSend(m.ChannelID, strings.Join(parts, "")+".png", bytes.NewReader(wttr))
		case "!wttrp2":
			var wttr []byte
			var err error
			if len(query) > 1 {
				wttr, err = weatherV2(url.QueryEscape(query[0]) + ".png" + "?" + query[1])
			} else {
				wttr, err = weatherV2(url.QueryEscape(query[0]) + ".png")
			}
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("wttr.in query failed")
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			s.ChannelFileSend(m.ChannelID, strings.Join(parts, "")+".png", bytes.NewReader(wttr))
		}

	}
}

func getWeather(baseURL string, location string) (result []byte, err error) {
	return httpGet(baseURL + location)
}

func weather(location string) (result []byte, err error) {
	return getWeather(baseURL, location)
}

func weatherV2(location string) (result []byte, err error) {
	return getWeather(baseURLv2, location)
}

func httpGet(url string) (result []byte, err error) {
	if _httpClient == nil {
		_httpClient = &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 300 * time.Second,
				}).Dial,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
	}

	var req *http.Request

	if req, err = http.NewRequest("GET", url, nil); err == nil {

		req.Header.Set("User-Agent", curlUserAgent)

		var resp *http.Response
		resp, err = _httpClient.Do(req)

		if resp != nil {
			defer resp.Body.Close() // in case of http redirects
		}

		if err == nil {
			var body []byte
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				return body, nil
			}
		}
	}

	return make([]byte, 0), err
}
