
$global/time

[1 2 3 $global/time]
(prn $global/time)

(= val 1)
$val

(if true [
    (prn 1)
] [
    (prn 2)
])


(= incr {
    [n] [
        (+ $n 1)
    ]
})

(map incr [1 2 3 4])
