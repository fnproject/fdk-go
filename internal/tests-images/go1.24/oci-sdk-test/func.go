/*
 * Copyright (c) 2021, 2022 Oracle and/or its affiliates. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/json"
	"io"
	"fmt"
    fdk "github.com/fnproject/fdk-go"
    "github.com/oracle/oci-go-sdk/v45/common/auth"
    "github.com/oracle/oci-go-sdk/v45/identity"
    "github.com/oracle/oci-go-sdk/v45/common"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(myHandler))
}

type Compart struct {
	CompartmentId string `json:"compartmentId"`
}

func myHandler(ctx context.Context, in io.Reader, out io.Writer) {
	    fmt.Println("Inside oci sdk test function")

        p := &Compart{CompartmentId: ""}
        json.NewDecoder(in).Decode(p)

        fmt.Println("compartmentId:", p.CompartmentId)
        if len(p.CompartmentId)==0 {
            panic("compartmentId cannot be empty/null")
        }
	    rp, err := auth.ResourcePrincipalConfigurationProvider()
    	if err != nil {
            fmt.Println(err)
            Msg:= fmt.Sprintf("Error", err)
            blob, err := json.Marshal(&Msg)
            if err != nil {
               panic(err)
            }
            written, err := out.Write(blob)
            if err != nil {
               panic(err)
            }
            if written != len(blob) {
               panic("Not all bytes written")
            }
           return
    	}
    	iam, err := identity.NewIdentityClientWithConfigurationProvider(rp)
       	if err != nil {
    		Msg:= fmt.Sprintf("Error", err)
            blob, err := json.Marshal(&Msg)
            if err != nil {
               panic(err)
            }
            written, err := out.Write(blob)
            if err != nil {
                panic(err)
            }
            if written != len(blob) {
                panic("Not all bytes written")
            }
            return
    	}
        req:= identity.GetCompartmentRequest{CompartmentId: common.String(p.CompartmentId)}
        resp, err := iam.GetCompartment(context.Background(), req)
        if err != nil {
            panic(err)
        }
        msg:= Compart{CompartmentId:*resp.Compartment.Id,}
        blob, err := json.Marshal(&msg)
        if err != nil {
            panic(err)
        }
        written, err := out.Write(blob)
        if err != nil {
            panic(err)
        }
        if written != len(blob) {
            panic("Not all bytes written")
       }
}
