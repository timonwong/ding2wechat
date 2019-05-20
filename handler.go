// Copyright (c) Timon Wong
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
	"github.com/timonwong/ding2wechat/config"
)

const (
	defaultTimeout = 10 * time.Second
)

type DingModel struct {
	MessageType string             `json:"msgtype"`
	Text        *DingModelText     `json:"text,omitemtpy"`
	Markdown    *DingModelMarkdown `json:"markdown,omitempty"`
}

type DingModelText struct {
	Content string `json:"content"`
}

type DingModelMarkdown struct {
	Text string `json:"text"`
}

type WechatModel struct {
	MessageType string               `json:"msgtype"`
	Text        *WechatModelText     `json:"text,omitemtpy"`
	Markdown    *WechatModelMarkdown `json:"markdown,omitempty"`
}

type WechatModelText struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list,omitempty"`
	MentionedMobileList []string `json:"mentioned_mobile_list,omitempty"`
}

type WechatModelMarkdown struct {
	Content string `json:"content"`
}

func translateModel(model *DingModel) (*WechatModel, error) {
	msgType := model.MessageType

	if msgType == "text" {
		if model.Text == nil {
			return nil, errors.New("malformed dingtalk body: missing text field")
		}

		return &WechatModel{
			MessageType: msgType,
			Text: &WechatModelText{
				Content: model.Text.Content,
			},
		}, nil
	} else if msgType == "markdown" {
		if model.Markdown == nil {
			return nil, errors.New("malformed dingtalk body: missing markdown field")
		}

		return &WechatModel{
			MessageType: msgType,
			Markdown: &WechatModelMarkdown{
				Content: model.Markdown.Text,
			},
		}, nil
	}

	return nil, fmt.Errorf("unknown msgtype: %s", msgType)
}

func ReceiverHandler(conf *config.Config) http.HandlerFunc {
	receivers := make(map[string]config.Receiver)
	for _, receiver := range conf.Receivers {
		receivers[receiver.Name] = receiver
	}

	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query()
		receiverName := qs.Get("name")

		receiver, ok := receivers[receiverName]
		if !ok {
			log.Errorf("Unknown receiver: %s", receiverName)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Parse incoming dingtalk model
		var dingModel DingModel
		err := json.NewDecoder(r.Body).Decode(&dingModel)
		if err != nil {
			log.Errorf("Unable to unmarshal dingtalk request body: %s", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		model, err := translateModel(&dingModel)
		if err != nil {
			log.Errorf("Unable to translate model: %s", err)
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		for _, target := range receiver.Targets {
			func() {
				modelCpy := *model
				mentionedList := target.MentionedList
				mentionedMobileList := target.MentionedMobileList

				if modelCpy.Text != nil {
					modelCpy.Text.MentionedList = mentionedList
					modelCpy.Text.MentionedMobileList = mentionedMobileList
				}

				b, err := json.Marshal(&modelCpy)
				if err != nil {
					log.Errorf("Unable to marshal model: %s", err)
					http.Error(w, "unknown error", http.StatusInternalServerError)
					return
				}

				log.Debugf("Sending http request to wechat webhook: %s", string(b))
				resp, err := httpClient.Post(target.URL, "application/json", bytes.NewReader(b))
				if err != nil {
					log.Errorf("Unable to send http request to wechat webhook: %s", err)
					http.Error(w, "http webhook error", http.StatusInternalServerError)
					return
				}
				defer resp.Body.Close()

				respBody, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					log.Debugf("Got response body: %s", string(respBody))
				}
			}()
		}

		w.Write([]byte("ok"))
	}
}
