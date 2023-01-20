import { parse } from "mfm-js"

globalThis["parse"] = (text: string) => JSON.stringify(parse(text))
