[[[ ]]]

$global/time


$~
$global/
$stack

(prn 5)

(set = $set)
(= val 3)
(prn $val)

(if (nil? $global/Name) [
    (set global/Name Ramen)
])

$global/Name

(prn
    Counter:
    $global/Counter
    ->
    (set global/Counter (+ $global/Counter 1))
)
