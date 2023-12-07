// Copyright 2021-2024 TimAndy. All rights reserved.
// Licensed under the Apache-2.0 license that can be found in the LICENSE file.

#include "funcdata.h"
#include "go_asm.h"
#include "textflag.h"

TEXT ·getgp(SB), NOSPLIT, $0-8
    MOVD    g, R8
    MOVD    R8, ret+0(FP)
    RET
