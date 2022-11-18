// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-4
    MOVW    g, R8
    MOVW    R8, ret+0(FP)
    RET

TEXT ·getg0(SB), NOSPLIT, $0-8
    NO_LOCAL_POINTERS
    MOVW    R8, ret_type+0(FP)
    MOVW    R9, ret_data+4(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOVW    $type·runtime·g(SB), R8
    //get runtime·g0 variable
    MOVW    $runtime·g0(SB), R9
    //return interface{}
    MOVW    R8, ret_type+0(FP)
    MOVW    R9, ret_data+4(FP)
    RET
