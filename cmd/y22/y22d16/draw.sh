#!/bin/bash

ls -1 *.dot | xargs -I{} dot -Tpng -o {}.png {}
