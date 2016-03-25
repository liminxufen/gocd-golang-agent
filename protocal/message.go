/*
 * Copyright 2016 ThoughtWorks, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package protocal

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

var (
	SetCookieAction           = "setCookie"
	CancelJobAction           = "cancelJob"
	ReregisterAction          = "reregister"
	BuildAction               = "build"
	PingAction                = "ping"
	AckAction                 = "ack"
	ReportCurrentStatusAction = "reportCurrentStatus"
	ReportCompletingAction    = "reportCompleting"
	ReportCompletedAction     = "reportCompleted"
)

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data"`
	AckId  string `json:"ackId"`
}

func (m *Message) Build() *Build {
	var build Build
	json.Unmarshal([]byte(m.Data), &build)
	return &build
}

func (m *Message) StringData() string {
	var str string
	json.Unmarshal([]byte(m.Data), &str)
	return str
}

func (m *Message) AgentRuntimeInfo() *AgentRuntimeInfo {
	var info AgentRuntimeInfo
	json.Unmarshal([]byte(m.Data), &info)
	return &info
}

func (m *Message) Report() *Report {
	var report Report
	json.Unmarshal([]byte(m.Data), &report)
	return &report
}

func newMessage(action string, data interface{}) *Message {
	json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return &Message{
		Action: action,
		Data:   string(json),
		AckId:  uuid.NewV4().String(),
	}
}

func SetCookieMessage(cookie string) *Message {
	return newMessage(SetCookieAction, cookie)
}

func AckMessage(ackId string) *Message {
	return newMessage(AckAction, ackId)
}

func BuildMessage(cmd *Build) *Message {
	return newMessage(BuildAction, cmd)
}

func PingMessage(data *AgentRuntimeInfo) *Message {
	return newMessage(PingAction, data)
}

func ReportMessage(t string, report *Report) *Message {
	return newMessage(t, report)
}

func CompletedMessage(report *Report) *Message {
	return ReportMessage(ReportCompletedAction, report)
}

func ReregisterMessage() *Message {
	return &Message{Action: ReregisterAction}
}
