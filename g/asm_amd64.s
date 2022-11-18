// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-8
    get_tls(CX)
    MOVQ    g(CX), AX
    MOVQ    AX, ret+0(FP)
    RET

TEXT ·getg0(SB), NOSPLIT, $0-16
    NO_LOCAL_POINTERS
    MOVQ    $0, ret_type+0(FP)
    MOVQ    $0, ret_data+8(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOVQ    $type·runtime·g(SB), AX
    //get runtime·g0 variable
    MOVQ    $runtime·g0(SB), BX
    //return interface{}
    MOVQ    AX, ret_type+0(FP)
    MOVQ    BX, ret_data+8(FP)
    RET
