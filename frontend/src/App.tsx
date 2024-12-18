import "./App.css";
import { BrowserRouter, Routes, Route } from "react-router-dom";

// Pages
import Layout from "./pages/main/Layout";
import Login from "./pages/login/Login";

function AppRouter() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout />} />
                <Route path="/login" element={<Login />} />
            </Routes>
        </BrowserRouter>
    );
}

function App() {
    return (
        <>
            <AppRouter />
        </>
    );
}

export default App;
