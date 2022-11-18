// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-8
    MOV    g, X10
    MOV    X10, ret+0(FP)
    RET

TEXT ·getg0(SB), NOSPLIT, $0-16
    NO_LOCAL_POINTERS
    MOV    X10, ret_type+0(FP)
    MOV    X11, ret_data+8(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOV    $type·runtime·g(SB), X10
    //get runtime·g0 variable
    MOV    $runtime·g0(SB), X11
    //return interface{}
    MOV    X10, ret_type+0(FP)
    MOV    X11, ret_data+8(FP)
    RET
