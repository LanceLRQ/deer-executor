<?php
while ($dat = fgets(STDIN)){
    $dat = explode(" ", $dat);
    if (count($dat) == 2) {
        echo sprintf("%d\n", intval($dat[0]) + intval($dat[1]));
    }
}