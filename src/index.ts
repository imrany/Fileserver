import express from "express"
import { config } from "dotenv"
import multer from "multer"
import cors from "cors"
import {readdir, mkdir} from "fs"
import fs from "fs"
import path from "path"
config()

const app =express()
const downloadPath = "./downloads";
if (!fs.existsSync(downloadPath)) {
    fs.mkdirSync(downloadPath, { recursive: true });
}

const storage = multer.memoryStorage();
const upload=multer({storage})
app.use(express.static("views"));
app.use(express.static("downloads"));
app.use(express.json({ limit: "10GB" }));
app.use(express.urlencoded({ limit: "10GB", extended: true }));
app.use(cors({}))

//routes
app.post("/upload", upload.array("files"), async(req:any, res:any) => {
    try {
        if (!req.files || !Array.isArray(req.files) || req.files.length === 0) {
            return res.status(400).json({ error: "No files received" });
        }
        const { fileName, chunkIndex, totalChunks } = req.body;
        if (!fileName) {
            return res.status(400).json({ error: "File name missing" });
        }
        const filePath = path.join(downloadPath, fileName);
        req.files.forEach((file:any) => {
            fs.appendFileSync(filePath, file.buffer);
        });
        console.log(`Uploaded chunk ${chunkIndex}/${totalChunks}`);
        res.status(200).json({ msg: `Chunk ${chunkIndex}/${totalChunks} received successfully!` });
    }
    catch (error:any) {
        console.error("Upload error:", error);
        res.status(500).json({ error: error.message });
    }
});

app.get("/read_file", async(_req:Express.Request, res:any) => {
    try {
        fs.readdir(downloadPath, "utf8", (err, files) => {
            if (err) {
                console.log(err);
                fs.mkdir(downloadPath, () => { console.log(`downloads dir made`); });
            }
            else {
                res.send(files);
            }
        });
    }
    catch (error:any) {
        res.status(505).send({ error: error.message });
    }
});

const port:any = process.env.PORT || 8000;
app.listen(port, "0.0.0.0", () => {
    console.log(`🚀 Server running at http://0.0.0.0:${port}`);
});
