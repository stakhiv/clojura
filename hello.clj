(def invert (fn (f)
  (fn (x)
    (not (f x)))))

(def even? (invert odd?))

(def x (range 10))
(println "odd?" (filter odd? x))
(println "even?" (filter even? x))

(println "odd sum" (reduce + 0 (filter odd? x)))
(println "even sum" (reduce + 0 (filter even? x)))
(println "total sum" (reduce + 0 x))
(println "map x2" (map (fn (x) (+ x x)) x))

(def blrd (fn (g)
  ;; (let score 0)
  (fn (n)
    ;; (set! score (inc score))
    (if (< n g)
      (println "More!")
      (if (> n g)
        (println "Less!")
        (println "You guessed it!!!"))))))

(def guess (blrd (random 100)))
