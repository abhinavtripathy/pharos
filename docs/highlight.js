// Colours the Pharos snippets on this site. The groups and their colours echo
// the highlighting in the project's original editor. Without JavaScript the
// code still reads fine, just in one colour.
(function () {
  "use strict";

  var GROUPS = [
    ["cmt", "--[^\\n]*"],
    ["str", "\"(?:[^\"\\\\\\n]|\\\\.)*\"|'(?:[^'\\\\\\n]|\\\\.)*'"],
    [
      "compare",
      "\\bis\\s+not\\s+equal\\s+to\\b|\\bis\\s+greater\\s+than\\b|\\bis\\s+less\\s+than\\b|\\bis\\s+equal\\s+to\\b"
    ],
    ["block", "\\bend\\s+codeblock\\b|\\bcodeblock\\b|\\bgive\\s+back\\b|\\bwith\\b"],
    [
      "logic",
      "\\botherwise\\s+if\\b|\\bend\\s+if\\b|\\botherwise\\b|\\bthen\\b|\\bif\\b|\\band\\b|\\bor\\b|\\bnot\\b"
    ],
    ["loop", "\\bend\\s+loop\\b|\\buntil\\b|\\bloop\\b"],
    ["io", "\\bprint\\b|\\boutput\\b|\\bout\\b"],
    ["type", "\\bnum\\b|\\bstring\\b|\\bbool\\b|\\btrue\\b|\\bfalse\\b"],
    ["num", "\\b\\d+(?:\\.\\d+)?\\b"]
  ];

  var PATTERN = new RegExp(
    GROUPS.map(function (g) {
      return "(" + g[1] + ")";
    }).join("|"),
    "gi"
  );

  function escapeHTML(text) {
    return text.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
  }

  function highlight(source) {
    return escapeHTML(source).replace(PATTERN, function () {
      // arguments[1..GROUPS.length] line up with the capture groups.
      for (var i = 0; i < GROUPS.length; i++) {
        var captured = arguments[i + 1];
        if (captured !== undefined) {
          return '<span class="tok-' + GROUPS[i][0] + '">' + captured + "</span>";
        }
      }
      return arguments[0];
    });
  }

  var blocks = document.querySelectorAll("code.pharos");
  for (var i = 0; i < blocks.length; i++) {
    blocks[i].innerHTML = highlight(blocks[i].textContent);
  }
})();
