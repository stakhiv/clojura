(def inc (fn
  (n)
  (+ n 1)))

(def dec (fn
  (n)
  (- n 1)))

(def each (fn (f l)
    ((if (head l)
       (do (f (head l))
           (recur f (tail l)))))))

; (def map (fn (f l)
;   (if l
;     (conj '((f (head l))) (map f (tail l))))))

; (def map2 (fn (f l)
;   (if l
;     (cons (f (head l)) (map f (tail l))))))

(def map (fn (f l)
  (_map3 f l '())))

(def _map3 (fn (f l r)
  (if l
    (recur f (tail l) (cons (f (head l)) r))
    r)))

; (def reduce (fn (f s l)
;     (if (head l)
;       (reduce f (f s (head l)) (tail l))
;       s)))

(def reduce (fn (f s l)
    (if (head l)
      (recur f (f s (head l)) (tail l))
      s)))

; (def range (fn (len)
;   (if (not (eq len 0))
;     (conj (range (dec len)) '(len)))))
;
; (def range2 (fn (len)
;   (if (not (eq len 0))
;     (cons len (range (dec len)))
;     '())))
;
; (def range3 (fn (len)
;   (_range3 len '())))
;
; (def _range3 (fn (len l)
;   (if (not (eq len 0))
;     (recur (dec len) (cons len l))
;     l)))
