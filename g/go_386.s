// Copyright 2022 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-4
    get_tls(CX)
    MOVL    g(CX), AX
    MOVL    AX, ret+0(FP)
    RET

TEXT ·getg0(SB), NOSPLIT, $0-8
    NO_LOCAL_POINTERS
    MOVL    $0, ret_type+0(FP)
    MOVL    $0, ret_data+4(FP)
    GO_RESULTS_INITIALIZED
    //get runtime.g type
    MOVL    $type·runtime·g(SB), AX
    //get runtime·g0 variable
    MOVL    $runtime·g0(SB), BX
    //return interface{}
    MOVL    AX, ret_type+0(FP)
    MOVL    BX, ret_data+4(FP)
    RET
