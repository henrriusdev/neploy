import { Icon, icons } from "@/lib/icons";
import { MagicWandIcon, PaperPlaneIcon } from "@radix-ui/react-icons";
import { Activity, AlarmClock, AlignCenter, Angry, Antenna, ArrowDown, ArrowLeft, ArrowRight, ArrowUp, AudioLines, AudioWaveform, Award, Axe, Bean, Beer, Bell, Binoculars, Blocks, Bomb, Box, Braces, Brain, Bug, Check, ChevronDown, ChevronLeft, ChevronRight, ChevronUp, Chrome, Code, CodeXml, Coffee, Database, Flame, FlaskConical, Folder, Gamepad2, Gem, Gift, Handshake, Headphones, Heart, Home, Hourglass, Key, Laptop, Lightbulb, Medal, Moon, Palette, Pencil, Phone, PiggyBank, Pin, Plane, Plug, Power, Rocket, Ruler, Scale, Scissors, Search, Shield, ShoppingCart, Skull, Snowflake, Speaker, Star, Sun, Sword, Tag, Target, Trophy, Umbrella, User, Users, Wallet, Wrench, X, Zap, ZoomIn, ZoomOut } from "lucide-react";

const iconComponent = (icon: Icon) => {
  switch (icon) {
  case "ArrowDown":
    return <ArrowDown className="w-10 h-10 text-white"/>;
  case "ArrowUp":
    return <ArrowUp className="w-10 h-10 text-white"/>;
  case "ArrowLeft":
    return <ArrowLeft className="w-10 h-10 text-white"/>;
  case "ArrowRight":
    return <ArrowRight className="w-10 h-10 text-white"/>;
  case "Check":
    return <Check className="w-10 h-10 text-white"/>;
  case "X":
    return <X className="w-10 h-10 text-white"/>;
  case "Search":
    return <Search className="w-10 h-10 text-white"/>;
  case "ChevronLeft":
    return <ChevronLeft className="w-10 h-10 text-white"/>;
  case "ChevronRight":
    return <ChevronRight className="w-10 h-10 text-white"/>;
  case "ChevronDown":
    return <ChevronDown className="w-10 h-10 text-white"/>;
  case "ChevronUp":
    return <ChevronUp className="w-10 h-10 text-white"/>;
  case "Activity":
    return <Activity className="w-10 h-10 text-white"/>;
  case "AlarmClock":
    return <AlarmClock className="w-10 h-10 text-white"/>;
  case "AlignCenter":
    return <AlignCenter className="w-10 h-10 text-white"/>;
  case "Angry":
    return <Angry className="w-10 h-10 text-white"/>;
  case "Antenna":
    return <Antenna className="w-10 h-10 text-white"/>;
  case "AudioLines":
    return <AudioLines className="w-10 h-10 text-white"/>;
  case "AudioWaveform":
    return <AudioWaveform className="w-10 h-10 text-white"/>;
  case "Award":
    return <Award className="w-10 h-10 text-white"/>;
  case "Axe":
    return <Axe className="w-10 h-10 text-white"/>;
  case "Bean":
    return <Bean className="w-10 h-10 text-white"/>;
  case "Beer":
    return <Beer className="w-10 h-10 text-white"/>;
  case "Bell":
    return <Bell className="w-10 h-10 text-white"/>;
  case "Binoculars":
    return <Binoculars className="w-10 h-10 text-white"/>;
  case "Blocks":
    return <Blocks className="w-10 h-10 text-white"/>;
  case "Bomb":
    return <Bomb className="w-10 h-10 text-white"/>;
  case "Box":
    return <Box className="w-10 h-10 text-white"/>;
  case "Braces":
    return <Braces className="w-10 h-10 text-white"/>;
  case "Brain":
    return <Brain className="w-10 h-10 text-white"/>;
  case "Bug":
    return <Bug className="w-10 h-10 text-white"/>;
  case "Chrome":
    return <Chrome className="w-10 h-10 text-white"/>;
  case "Code":
    return <Code className="w-10 h-10 text-white"/>;
  case "CodeXml":
    return <CodeXml className="w-10 h-10 text-white"/>;
  case "Coffee":
    return <Coffee className="w-10 h-10 text-white"/>;
  case "Database":
    return <Database className="w-10 h-10 text-white"/>;
  case "FlaskConical":
    return <FlaskConical className="w-10 h-10 text-white"/>;
  case "Flame":
    return <Flame className="w-10 h-10 text-white"/>;
  case "Folder":
    return <Folder className="w-10 h-10 text-white"/>;
  case "Gamepad2":
    return <Gamepad2 className="w-10 h-10 text-white"/>;
  case "Gem":
    return <Gem className="w-10 h-10 text-white"/>;
  case "Gift":
    return <Gift className="w-10 h-10 text-white"/>;
  case "Handshake":
    return <Handshake className="w-10 h-10 text-white"/>;
  case "Headphones":
    return <Headphones className="w-10 h-10 text-white"/>;
  case "Heart":
    return <Heart className="w-10 h-10 text-white"/>;
  case "Home":
    return <Home className="w-10 h-10 text-white"/>;
  case "Hourglass":
    return <Hourglass className="w-10 h-10 text-white"/>;
  case "Key":
    return <Key className="w-10 h-10 text-white"/>;
  case "Laptop":
    return <Laptop className="w-10 h-10 text-white"/>;
  case "Lightbulb":
    return <Lightbulb className="w-10 h-10 text-white"/>;
  case "Lock":
    return <Lock className="w-10 h-10 text-white"/>;
  case "MagicWand":
    return <MagicWandIcon className="w-10 h-10 text-white"/>;
  case "Map":
    return <Map className="w-10 h-10 text-white"/>;
  case "Medal":
    return <Medal className="w-10 h-10 text-white"/>;
  case "Moon":
    return <Moon className="w-10 h-10 text-white"/>;
  case "Palette":
    return <Palette className="w-10 h-10 text-white"/>;
  case "PaperPlane":
    return <PaperPlaneIcon className="w-10 h-10 text-white"/>;
  case "Pencil":
    return <Pencil className="w-10 h-10 text-white"/>;
  case "Phone":
    return <Phone className="w-10 h-10 text-white"/>;
  case "PiggyBank":
    return <PiggyBank className="w-10 h-10 text-white"/>;
  case "Pin":
    return <Pin className="w-10 h-10 text-white"/>;
  case "Plane":
    return <Plane className="w-10 h-10 text-white"/>;
  case "Plug":
    return <Plug className="w-10 h-10 text-white"/>;
  case "Power":
    return <Power className="w-10 h-10 text-white"/>;
  case "Rocket":
    return <Rocket className="w-10 h-10 text-white"/>;
  case "Ruler":
    return <Ruler className="w-10 h-10 text-white"/>;
  case "Scale":
    return <Scale className="w-10 h-10 text-white"/>;
  case "Scissors":
    return <Scissors className="w-10 h-10 text-white"/>;
  case "Shield":
    return <Shield className="w-10 h-10 text-white"/>;
  case "ShoppingCart":
    return <ShoppingCart className="w-10 h-10 text-white"/>;
  case "Skull":
    return <Skull className="w-10 h-10 text-white"/>;
  case "Snowflake":
    return <Snowflake className="w-10 h-10 text-white"/>;
  case "Speaker":
    return <Speaker className="w-10 h-10 text-white"/>;
  case "Star":
    return <Star className="w-10 h-10 text-white"/>;
  case "Sun":
    return <Sun className="w-10 h-10 text-white"/>;
  case "Sword":
    return <Sword className="w-10 h-10 text-white"/>;
  case "Tag":
    return <Tag className="w-10 h-10 text-white"/>;
  case "Target":
    return <Target className="w-10 h-10 text-white"/>;
  case "Trophy":
    return <Trophy className="w-10 h-10 text-white"/>;
  case "Umbrella":
    return <Umbrella className="w-10 h-10 text-white"/>;
  case "User":
    return <User className="w-10 h-10 text-white"/>;
  case "Users":
    return <Users className="w-10 h-10 text-white"/>;
  case "Wallet":
    return <Wallet className="w-10 h-10 text-white"/>;
  case "Wrench":
    return <Wrench className="w-10 h-10 text-white"/>;
  case "ZoomIn":
    return <ZoomIn className="w-10 h-10 text-white"/>;
  case "ZoomOut":
    return <ZoomOut className="w-10 h-10 text-white"/>;
  case "Zap":
    return <Zap className="w-10 h-10 text-white"/>;
}
}
export const RoleIcon = ({ icon, color }: { icon: string; color: string }) => {
  const iconTyped = icon as Icon;
  return (
    <div
      className="flex items-center justify-center w-14 h-14 rounded-lg"
      style={{ backgroundColor: color }}
    >
      {iconComponent(iconTyped)}
    </div>
  );
};
