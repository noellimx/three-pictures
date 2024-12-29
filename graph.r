# Hello, world!
#
# This is an example function named 'hello'
# which prints 'Hello, world!'.
#
# You can learn more about package authoring with RStudio at:
#
#   https://r-pkgs.org
#
# Some useful keyboard shortcuts for package authoring:
#
#   Install Package:           'Cmd + Shift + B'
#   Check Package:             'Cmd + Shift + E'
#   Test Package:              'Cmd + Shift + T'
library(scatterplot3d)

hello <- function() {
  print("Hello, world!")
}


outcome_data <- read.csv("./some.csv")


head(outcome_data)

outcome_data <- outcome_data[order(outcome_data$average_player_pct_change), ]

outcome_data$half_pay_on[outcome_data$half_pay_on == -1] <- 11
X <- as.numeric(unlist(outcome_data["rounds_played"]))
Y <- as.numeric(unlist(outcome_data$half_pay_on))
# Y[Y == -1] <- 11
Z <- as.numeric(unlist(outcome_data["average_player_pct_change"]))

head(Y)
head(outcome_data)

colors <- c(
  "#FF5733", "#FF5733", "#FF5733", "#FF5733", "#FF5733", "#FF5733",

  #"#FF5733", "#33FF57", "#3357FF", "#FF33A8", "#A833FF", "#33FFF1",
  # "#FFF133", "#F133FF", "#33A8FF", "#57FF33", "#5733FF", "#FF3357", "#A8FF33"
  "blue", "blue", "blue", "black", "black", "black", "black"
)
colll <- colors[unlist(Y)]
colll <- ifelse(unlist(Y) == 0, "orange", colors[unlist(Y)])


length(Y)
length(unlist(Y))

length(X)
length(colors)
length(colll)
show(colll)
show(Y)

p_color <- seq(1,7)
p_color

shapes <- p_color[outcome_data$players]

for (i in seq(25, 200, by = 5)) {
  scatterplot3d(main="Buyin = 1000Blind",x=X,xlab = "rounds_played",y =Y,ylab = "half_pay_on",z=Z,zlab = "percentage change",
                angle = i,color=colll,pch=shapes)
}


# Create a color gradient based on the x values
colors <- colorRampPalette(c("lightblue","blue"))(100)  # Define a gradient from blue to red

plot(main="Buyin = 1000Blind",x =Y,xlab = "half_pay_on",y=Z,ylab = "percentage change",col = colors[as.integer(cut(X, breaks = 100))], pch = 16)
