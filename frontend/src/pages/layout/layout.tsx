import { useState, useEffect } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
    Home,
    Settings,
    Users,
    FileText,
    Menu,
    CalendarCheck,
    Package,
    ChevronRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const Layout = () => {
    const [userRole, setUserRole] = useState(0);
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [currentTab, setCurrentTab] = useState("dashboard");
    const [currentSubTab, setCurrentSubTab] = useState("");
    const [windowWidth, setWindowWidth] = useState(
        typeof window !== "undefined" ? window.innerWidth : 0
    );

    useEffect(() => {
        setUserRole(3); // for dev - setting as admin
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
            label: "Lab Books",
            id: "labbooks",
            role: 1,
            subTabs: [
                { name: "Upload", role: 1 },
                { name: "Approval", role: 2 },
                { name: "History", role: 1 },
            ],
        },
        {
            icon: <CalendarCheck className="w-4 h-4" />,
            label: "Schedules",
            id: "schedules",
            role: 1,
            subTabs: [],
        },
        {
            icon: <Package className="w-4 h-4" />,
            label: "Resources",
            id: "resources",
            role: 1,
            subTabs: [],
        },
        {
            icon: <Users className="w-4 h-4" />,
            label: "Users",
            id: "users",
            role: 2,
            subTabs: [
                { name: "All Users", role: 2 },
                { name: "Register", role: 2 },
                { name: "Permissions", role: 3 },
            ],
        },
        {
            icon: <Settings className="w-4 h-4" />,
            label: "Settings",
            id: "settings",
            role: 1,
            subTabs: [
                { name: "General", role: 1 },
                { name: "Security", role: 1 },
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

    const Navigation = () => {
        const handleItemClick = (itemId: string) => {
            setCurrentTab(itemId);
            setCurrentSubTab("");
            if (!isDesktop) {
                setIsMobileOpen(false);
            }
        };

        const handleSubTabClick = (itemId: string, subTabName: string) => {
            setCurrentTab(itemId);
            setCurrentSubTab(subTabName);
            if (!isDesktop) {
                setIsMobileOpen(false);
            }
        };

        return (
            <div className="h-full flex flex-col">
                <div className="h-14 flex flex-col items-center justify-center px-4 border-b">
                    <h2 className="text-lg font-bold">AlphaLab</h2>
                    <span className="text-xs text-gray-500">v0.0.1</span>
                </div>

                <ScrollArea className="flex-1 py-2">
                    <nav className="space-y-1 px-2">
                        {navItems
                            .filter((item) => userRole >= item.role)
                            .map((item) => {
                                const hasSubTabs = item.subTabs.some(
                                    (subTab) => userRole >= subTab.role
                                );

                                const filteredSubTabs = item.subTabs.filter(
                                    (subTab) => userRole >= subTab.role
                                );

                                if (!hasSubTabs) {
                                    return (
                                        <div
                                            key={item.id}
                                            role="button"
                                            tabIndex={0}
                                            onClick={() =>
                                                handleItemClick(item.id)
                                            }
                                            className={cn(
                                                "w-full flex items-center gap-3 rounded-lg px-3 py-2 transition-all duration-200 cursor-pointer hover:text-gray-900 hover:bg-gray-100",
                                                currentTab === item.id &&
                                                    currentSubTab === ""
                                                    ? "bg-gray-100 text-gray-900"
                                                    : "text-gray-500"
                                            )}
                                        >
                                            {item.icon}
                                            <span>{item.label}</span>
                                        </div>
                                    );
                                }

                                return (
                                    <DropdownMenu key={item.id}>
                                        <DropdownMenuTrigger asChild>
                                            <div
                                                className={cn(
                                                    "w-full flex items-center justify-between rounded-lg px-3 py-2 transition-all duration-200 cursor-pointer hover:text-gray-900 hover:bg-gray-100",
                                                    currentTab === item.id
                                                        ? "bg-gray-100 text-gray-900"
                                                        : "text-gray-500"
                                                )}
                                            >
                                                <div className="flex items-center gap-3">
                                                    {item.icon}
                                                    <span>{item.label}</span>
                                                </div>
                                                <ChevronRight className="h-4 w-4" />
                                            </div>
                                        </DropdownMenuTrigger>
                                        <DropdownMenuContent
                                            side="right"
                                            align="start"
                                            className="w-48"
                                        >
                                            {filteredSubTabs.map((subTab) => (
                                                <DropdownMenuItem
                                                    key={subTab.name}
                                                    onClick={() =>
                                                        handleSubTabClick(
                                                            item.id,
                                                            subTab.name
                                                        )
                                                    }
                                                    className={cn(
                                                        "cursor-pointer",
                                                        currentTab ===
                                                            item.id &&
                                                            currentSubTab ===
                                                                subTab.name
                                                            ? "bg-gray-100"
                                                            : ""
                                                    )}
                                                >
                                                    {subTab.name}
                                                </DropdownMenuItem>
                                            ))}
                                        </DropdownMenuContent>
                                    </DropdownMenu>
                                );
                            })}
                    </nav>
                </ScrollArea>
            </div>
        );
    };

    // Get current navigation path for title
    const getNavigationPath = () => {
        const mainTab = getCurrentTabLabel();

        return currentSubTab ? (
            <span className="flex items-center">
                <span className="text-gray-500">{mainTab}</span>
                <span className="mx-2 text-gray-500">/</span>
                <span className="text-black">{currentSubTab}</span>
            </span>
        ) : (
            <span className="text-gray-500">{mainTab}</span>
        );
    };

    const getCurrentTabLabel = () => {
        return navItems.find((item) => item.id === currentTab)?.label || "";
    };

    return (
        <div className="min-h-screen bg-background">
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
                    {/* Navigation Path Title */}
                    <div className="h-14 flex items-center px-6">
                        <div className="flex items-center">
                            {!isDesktop && <div className="w-8" />}
                            <h1 className="text-xl font-semibold">
                                {getNavigationPath()}
                            </h1>
                        </div>
                    </div>
                </header>

                {/* Main Content Area */}
                <div className="flex-1 flex flex-col overflow-auto">
                    <main className="flex-1 p-6">
                        <div className="rounded-lg border p-4">
                            <p>
                                Content for{" "}
                                {currentSubTab || getCurrentTabLabel()} will go
                                here
                            </p>
                        </div>
                    </main>

                    <footer className="flex-none border-t p-3 text-center bg-background">
                        <p className="text-[#B4B4B4]">
                            Copyright ©{" "}
                            <a
                                href="https://github.com/JRay-Lin/AlphaLab"
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
