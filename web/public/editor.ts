import { defaultValueCtx, Editor } from "@milkdown/core";
import { nord } from "@milkdown/theme-nord";
import { commonmark } from "@milkdown/preset-commonmark";
import { emoji } from "@milkdown/plugin-emoji";


Editor.make()
    .config((ctx) => {
        ctx.set(defaultValueCtx, "# Milkdown :heartpulse: Vanilla");
    })
    .use(nord)
    .use(commonmark)
    .use(emoji)
    .create();
