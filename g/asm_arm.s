// Copyright 2021-2025 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getg(SB), NOSPLIT, $0-4
    MOVW    g, R8
    MOVW    R8, ret+0(FP)
    RET
