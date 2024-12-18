import { useState, useEffect } from "react";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";
import { Navigation, navItems, getComponent } from "./Nav";

const Layout = () => {
    const [userRole, setUserRole] = useState(0);
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [currentTab, setCurrentTab] = useState("dashboard");
    const [currentSubTab, setCurrentSubTab] = useState("");
    const [windowWidth, setWindowWidth] = useState(
        typeof window !== "undefined" ? window.innerWidth : 0
    );

    useEffect(() => {
        setUserRole(3); // For dev - setting as admin
    }, []);

    useEffect(() => {
        const handleResize = () => {
            setWindowWidth(window.innerWidth);
        };

        window.addEventListener("resize", handleResize);
        return () => window.removeEventListener("resize", handleResize);
    }, []);

    useEffect(() => {
        const currentNavItem = navItems.find((item) => item.id === currentTab);
        if (currentNavItem) {
            const availableSubTabs = currentNavItem.subTabs.filter(
                (subTab) => userRole >= subTab.role
            );
            if (availableSubTabs.length > 0) {
                setCurrentSubTab(availableSubTabs[0].name);
            } else {
                // Clear subtab if the current tab doesn't have any subtabs
                setCurrentSubTab("");
            }
        }
    }, [currentTab, userRole]);

    const isDesktop = windowWidth >= 1024;

    const getNavigationPath = () => {
        const mainTab =
            navItems.find((item) => item.id === currentTab)?.label || "";
        const currentNavItem = navItems.find((item) => item.id === currentTab);
        const hasSubTabs = currentNavItem?.subTabs.some(
            (subTab) => userRole >= subTab.role
        );

        return hasSubTabs && currentSubTab ? (
            <span className="flex items-center">
                <span className="text-gray-500">{mainTab}</span>
                <span className="mx-2 text-gray-500">/</span>
                <span className="text-black">{currentSubTab}</span>
            </span>
        ) : (
            <span className="text-gray-500">{mainTab}</span>
        );
    };

    const handleTabChange = (tabId: string) => {
        setCurrentTab(tabId);
        const newNavItem = navItems.find((item) => item.id === tabId);
        if (!newNavItem?.subTabs.length) {
            setCurrentSubTab("");
        }
    };

    const renderContent = () => {
        const Component = getComponent(currentTab, currentSubTab);

        if (Component) {
            return <Component />;
        }

        return (
            <div className="rounded-lg border p-4">
                <p>
                    Content for{" "}
                    {currentSubTab ||
                        navItems.find((item) => item.id === currentTab)
                            ?.label}{" "}
                    is not yet implemented.
                </p>
            </div>
        );
    };

    return (
        <div className="min-h-screen bg-background relative">
            {isDesktop && (
                <div
                    className="fixed left-0 top-0 h-full w-64 border-r bg-background"
                    style={{ zIndex: 40 }}
                >
                    <Navigation
                        currentTab={currentTab}
                        currentSubTab={currentSubTab}
                        userRole={userRole}
                        onItemClick={handleTabChange}
                        onSubTabClick={(itemId, subTabName) => {
                            setCurrentTab(itemId);
                            setCurrentSubTab(subTabName);
                        }}
                        isDesktop={isDesktop}
                    />
                </div>
            )}
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
                        <Navigation
                            currentTab={currentTab}
                            currentSubTab={currentSubTab}
                            userRole={userRole}
                            onItemClick={handleTabChange}
                            onSubTabClick={(itemId, subTabName) => {
                                setCurrentTab(itemId);
                                setCurrentSubTab(subTabName);
                            }}
                            isDesktop={isDesktop}
                            onMobileClose={() => setIsMobileOpen(false)}
                        />
                    </SheetContent>
                </Sheet>
            )}
            <div
                className={`${
                    isDesktop ? "lg:pl-64" : ""
                } flex flex-col h-screen relative`}
                style={{ zIndex: 30 }}
            >
                <header className="flex-none border-b bg-background z-20">
                    <div className="h-14 flex items-center px-6">
                        <h1 className="text-xl font-semibold">
                            {getNavigationPath()}
                        </h1>
                    </div>
                </header>
                <div className="flex-1 flex flex-col overflow-auto relative z-10">
                    <main className="flex-1 p-6">{renderContent()}</main>
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
