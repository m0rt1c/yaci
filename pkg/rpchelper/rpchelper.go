package rpchelper

import (
  "../chord"
)

type (
  // ServiceArgs arguments for yaci service rpc
  ServiceArgs struct {
    Name string
    Port int
    Key string
    Base int
    Exponent int
  }

  // ServiceReply reply for yaci service rpc
  ServiceReply struct {
    Node chord.NodeInfo
    Message string
    List map[string]chord.Node
  }
)