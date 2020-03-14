package rpc

/*
	Sliver Implant Framework
	Copyright (C) 2019  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"context"
	"time"

	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/protobuf/sliverpb"
	"github.com/bishopfox/sliver/server/core"
	"github.com/golang/protobuf/proto"
)

// GetSessions - Get a list of sessions
func (rpc *Server) GetSessions(ctx context.Context, _ *commonpb.Empty) (*clientpb.Sessions, error) {
	sessions := &clientpb.Sessions{
		Sessions: []*clientpb.Session{},
	}
	if 0 < len(*core.Hive.Slivers) {
		for _, sliver := range *core.Hive.Slivers {
			sessions.Sessions = append(sessions.Sessions, sliver.ToProtobuf())
		}
	}
	return sessions, nil
}

// KillSession - Kill a session
func (rpc *Server) KillSession(ctx context.Context, kill *sliverpb.KillSessionReq) (*commonpb.Empty, error) {
	sliver := core.Hive.Sliver(kill.Request.SessionID)
	if sliver == nil {
		return &commonpb.Empty{}, ErrInvalidSessionID
	}
	core.Hive.RemoveSliver(sliver)
	data, err := proto.Marshal(kill)
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(kill.Request.GetTimeout())
	sliver.Request(sliverpb.MsgNumber(kill), timeout, data)
	return &commonpb.Empty{}, nil
}
