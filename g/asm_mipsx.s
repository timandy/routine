// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

//go:build mips || mipsle
// +build mips mipsle

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-4
    MOVW    g, R8
    MOVW    R8, ret+0(FP)
    RET
