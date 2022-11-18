// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-8
    MOVD    g, R8
    MOVD    R8, ret+0(FP)
    RET

TEXT ·getg0(SB), NOSPLIT, $0-16
    NO_LOCAL_POINTERS
    MOVD    $0, ret_type+0(FP)
    MOVD    $0, ret_data+8(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOVD    $type·runtime·g(SB), R8
    //get runtime·g0 variable
    MOVD    $runtime·g0(SB), R9
    //return interface{}
    MOVD    R8, ret_type+0(FP)
    MOVD    R9, ret_data+8(FP)
    RET
