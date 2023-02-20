// Copyright 2023 TimAndy. All rights reserved.
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
