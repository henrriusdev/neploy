import { Check } from "lucide-react";


 export const RenderStepIndicators = ({ step, totalSteps, steps }) => (
   <div className="flex justify-center mb-8">
     {steps.map((_, index) => (
       <div key={index} className="flex items-center">
         <div
           className={`w-8 h-8 rounded-full flex items-center justify-center ${
             step > index + 1
               ? "bg-lime-700 text-white"
               : step === index + 1
               ? "bg-primary text-white"
               : "bg-gray-200"
           }`}>
           {step > index + 1 ? <Check className="w-4 h-4" /> : index + 1}
         </div>
         {index < totalSteps - 1 && (
           <div
             className={`w-6 sm:w-8 md:w-10 lg:w-12 h-1 ${
               step - 1 === index + 1
                 ? "bg-gradient-to-r from-lime-700 to-primary from-40% to-90%"
                 : step > index + 1
                 ? "bg-lime-700"
                 : "bg-gray-200"
             }`}
           />
         )}
       </div>
     ))}
   </div>
 );