// Copyright 2023 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

//go:build mips64 || mips64le
// +build mips64 mips64le

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-8
    MOVV    g, R8
    MOVV    R8, ret+0(FP)
    RET
