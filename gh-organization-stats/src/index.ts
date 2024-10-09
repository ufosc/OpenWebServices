import express from "express";
import mongoose from "mongoose";
import clubRoutes from "./routes";
import dotenv from "dotenv";
dotenv.config();

const app = express();
app.use(express.json());

const mongoUri = process.env.DB_KEY;
if (!mongoUri) {
  throw new Error(
    "MongoDB connection string is missing in environment variables.",
  );
}

mongoose
  .connect(mongoUri)
  .then(() => {
    console.log("MongoDB connected");
  })
  .catch((err) => console.error("Error connecting to MongoDB: ", err));

app.use(express.json());
app.use("/api", clubRoutes);

app.get("/", (req, res) => {
  res.send("API is running...");
});

const PORT = process.env.PORT || 8000;
app.listen(PORT, () => console.log(`Server running on port ${PORT}`));
