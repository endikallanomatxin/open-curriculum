import test from "node:test";
import assert from "node:assert/strict";
import {
  safeFileName,
  slugify,
  topologicalLevels,
  wouldCreateCycle,
} from "../worker/lib.js";

test("slugify normalizes accents and punctuation", () => {
  assert.equal(slugify("Álgebra lineal: vectores"), "algebra-lineal-vectores");
});

test("safeFileName strips unsafe path characters", () => {
  assert.equal(safeFileName("../../Guía final (v2).pdf"), "..-..-Guia-final-v2-.pdf");
});

test("cycle detection rejects direct and indirect cycles", () => {
  const edges = [
    { unit_id: "b", prerequisite_id: "a" },
    { unit_id: "c", prerequisite_id: "b" },
  ];
  assert.equal(wouldCreateCycle(edges, "a", "c"), true);
  assert.equal(wouldCreateCycle(edges, "d", "c"), false);
  assert.equal(wouldCreateCycle(edges, "a", "a"), true);
});

test("topologicalLevels places prerequisites first", () => {
  const units = [{ id: "a" }, { id: "b" }, { id: "c" }];
  const levels = topologicalLevels(units, [
    { unit_id: "b", prerequisite_id: "a" },
    { unit_id: "c", prerequisite_id: "b" },
  ]);
  assert.deepEqual(levels, { a: 0, b: 1, c: 2 });
});
