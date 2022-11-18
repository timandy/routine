// Copyright 2022 TimAndy. All rights reserved.
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

TEXT ·getg0(SB), NOSPLIT, $0-16
    NO_LOCAL_POINTERS
    MOVV    $0, ret_type+0(FP)
    MOVV    $0, ret_data+8(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOVV    $type·runtime·g(SB), R8
    //get runtime·g0 variable
    MOVV    $runtime·g0(SB), R9
    //return interface{}
    MOVV    R8, ret_type+0(FP)
    MOVV    R9, ret_data+8(FP)
    RET
