import React, { useState, useEffect } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    Home,
    Settings,
    Users,
    FileText,
    Menu,
    CalendarCheck,
    Package,
} from "lucide-react";
import { Button } from "@/components/ui/button";

const Layout = () => {
    const [userRole, setUserRole] = useState(0);
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [currentTab, setCurrentTab] = useState("dashboard");
    const [currentSubTab, setCurrentSubTab] = useState("");
    const [windowWidth, setWindowWidth] = useState(
        typeof window !== "undefined" ? window.innerWidth : 0
    );

    useEffect(() => {
        setUserRole(1); // for dev - setting as admin
    }, []);

    const navItems = [
        {
            icon: <Home className="w-4 h-4" />,
            label: "Dashboard",
            id: "dashboard",
            role: 0,
            subTabs: [
                { name: "Overview", role: 0 },
                { name: "Analytics", role: 1 },
                { name: "Reports", role: 2 },
                { name: "Metrics", role: 2 },
            ],
        },
        {
            icon: <FileText className="w-4 h-4" />,
            label: "Reports",
            id: "reports",
            role: 1,
            subTabs: [
                { name: "Upload", role: 1 },
                { name: "Approval", role: 2 },
                { name: "History", role: 1 },
                { name: "Archive", role: 2 },
            ],
        },
        {
            icon: <CalendarCheck className="w-4 h-4" />,
            label: "Schedule",
            id: "schedules",
            role: 1,
            subTabs: [
                { name: "List", role: 1 },
                { name: "Create", role: 2 },
            ],
        },
        {
            icon: <Package className="w-4 h-4" />,
            label: "Resources",
            id: "resources",
            role: 1,
            subTabs: [
                { name: "List", role: 1 },
                { name: "Tags", role: 2 },
            ],
        },
        {
            icon: <Users className="w-4 h-4" />,
            label: "Users",
            id: "users",
            role: 2,
            subTabs: [
                { name: "All Users", role: 2 },
                { name: "Permissions", role: 3 },
                { name: "Register", role: 3 },
            ],
        },
        {
            icon: <Settings className="w-4 h-4" />,
            label: "Settings",
            id: "settings",
            role: 3,
            subTabs: [
                { name: "General", role: 3 },
                { name: "Security", role: 3 },
                { name: "Notifications", role: 3 },
                { name: "Integrations", role: 3 },
            ],
        },
    ];

    useEffect(() => {
        let timeoutId: NodeJS.Timeout;

        const handleResize = () => {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => {
                setWindowWidth(window.innerWidth);
            }, 100);
        };

        window.addEventListener("resize", handleResize);
        return () => {
            window.removeEventListener("resize", handleResize);
            clearTimeout(timeoutId);
        };
    }, []);

    // Set initial subtab when main tab changes
    useEffect(() => {
        const currentNavItem = navItems.find((item) => item.id === currentTab);
        if (currentNavItem) {
            const availableSubTabs = currentNavItem.subTabs.filter(
                (subTab) => userRole >= subTab.role
            );
            if (availableSubTabs.length > 0) {
                setCurrentSubTab(availableSubTabs[0].name);
            }
        }
    }, [currentTab, userRole]);

    const isDesktop = windowWidth >= 1024;

    const Navigation = () => (
        <div className="h-full flex flex-col">
            <div className="h-14 flex items-center px-4 border-b">
                <h2 className="text-lg font-semibold">My App</h2>
            </div>

            <ScrollArea className="flex-1 py-2">
                <nav className="space-y-1 px-2">
                    {navItems
                        .filter((item) => userRole >= item.role)
                        .map((item) => (
                            <a
                                key={item.id}
                                href={`#${item.id}`}
                                onClick={(e) => {
                                    e.preventDefault();
                                    setCurrentTab(item.id);
                                    if (!isDesktop) {
                                        setIsMobileOpen(false);
                                    }
                                }}
                                className={`flex items-center gap-3 rounded-lg px-3 py-2 transition-all hover:text-gray-900 hover:bg-gray-100 ${
                                    currentTab === item.id
                                        ? "bg-gray-100 text-gray-900"
                                        : "text-gray-500"
                                }`}
                            >
                                {item.icon}
                                <span>{item.label}</span>
                            </a>
                        ))}
                </nav>
            </ScrollArea>
        </div>
    );

    const getCurrentTabLabel = () => {
        return navItems.find((item) => item.id === currentTab)?.label || "";
    };

    const getCurrentSubTabs = () => {
        const currentNavItem = navItems.find((item) => item.id === currentTab);
        if (!currentNavItem) return [];
        return currentNavItem.subTabs
            .filter((subTab) => userRole >= subTab.role)
            .map((subTab) => subTab.name);
    };

    return (
        <div className="min-h-screen bg-background">
            {/* Rest of the layout code remains the same... */}
            {/* Desktop Navigation */}
            {isDesktop && (
                <div className="fixed left-0 top-0 h-full w-64 border-r bg-background z-30">
                    <Navigation />
                </div>
            )}

            {/* Mobile Navigation */}
            {!isDesktop && (
                <Sheet open={isMobileOpen} onOpenChange={setIsMobileOpen}>
                    <SheetTrigger asChild>
                        <Button
                            variant="ghost"
                            size="icon"
                            className="fixed top-3 left-3 z-50"
                        >
                            <Menu className="h-6 w-6" />
                        </Button>
                    </SheetTrigger>
                    <SheetContent side="left" className="w-64 p-0">
                        <Navigation />
                    </SheetContent>
                </Sheet>
            )}

            {/* Main Content */}
            <div
                className={`${
                    isDesktop ? "lg:pl-64" : ""
                } flex flex-col h-screen`}
            >
                {/* Fixed Header */}
                <header className="flex-none border-b bg-background z-20">
                    {/* Current Section Title */}
                    <div className="border-b h-14 flex items-center px-6">
                        <div className="flex items-center">
                            {!isDesktop && <div className="w-8" />}
                            <h1 className="text-xl font-semibold">
                                {getCurrentTabLabel()}
                            </h1>
                        </div>
                    </div>

                    {/* Sub Features Tabs */}
                    <Tabs
                        value={currentSubTab}
                        onValueChange={setCurrentSubTab}
                        className="w-full"
                    >
                        <TabsList className="w-full justify-start rounded-none border-b h-12 px-6">
                            {!isDesktop && <div className="w-12" />}
                            <div className="flex-1 flex gap-2">
                                {getCurrentSubTabs().map((subTab) => (
                                    <TabsTrigger
                                        key={subTab}
                                        value={subTab}
                                        className="data-[state=active]:bg-background px-4"
                                    >
                                        {subTab}
                                    </TabsTrigger>
                                ))}
                            </div>
                        </TabsList>
                    </Tabs>
                </header>

                {/* Scrollable Content Area */}
                <div className="flex-1 flex flex-col overflow-auto">
                    {/* Main Tab Content */}
                    <main className="flex-1 p-6">
                        <div className="rounded-lg border p-4">
                            <h2 className="text-lg font-medium mb-2">
                                {getCurrentTabLabel()} - {currentSubTab}
                            </h2>
                            <p>Content for {currentSubTab} will go here</p>
                        </div>
                    </main>

                    {/* Footer */}
                    <footer className="flex-none border-t p-4 text-center bg-background">
                        <p className="text-[#B4B4B4]">
                            Copyright ©{" "}
                            <a
                                href="https://github.com/JRay9487"
                                className="text-[#4096B6] hover:underline"
                                target="_blank"
                                rel="noopener noreferrer"
                            >
                                AlphaLab
                            </a>{" "}
                            {new Date().getFullYear()}. All Rights Reserved.
                        </p>
                    </footer>
                </div>
            </div>
        </div>
    );
};

export default Layout;
